package cif

import (
  bolt "github.com/coreos/bbolt"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/statistics"
  "log"
  "net/http"
  "time"
)

// GetSchedule returns a specific Schedule's for a specific TrainUID, startDate and STPIndicator
// If no schedule exists for the required key then nil is returned
func (c *CIF) GetSchedule( tx *bolt.Tx, uid string, date time.Time, stp string ) *Schedule {
  key := []byte( uid + date.Format( Date ) + stp )
  s := &Schedule{}
  b := tx.Bucket( []byte( "Schedule" ) ).Get( key )
  dec := codec.NewBinaryCodecFrom( b )
  dec.Read( s )
  if uid == s.TrainUID {
    return s
  }
  return nil
}

// ScheduleHandler implements a net/http handler that implements a simple Rest
// service to retrieve all schedules for a specific uid, date and STPIndicator
// The handler must have {uid} set in the path for this to work.
//
// For example:
//
// router.HandleFunc( "/schedule/{uid}/{date}/{stp}", db.ScheduleHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct.
func (c *CIF) ScheduleHandler( w http.ResponseWriter, r *http.Request ) {
  var params = mux.Vars( r )

  key := params[ "uid" ] + params[ "date" ] + params[ "stp" ]

  if err := c.db.View( func( tx *bolt.Tx ) error {

    s := &Schedule{}
    b := tx.Bucket( []byte( "Schedule" ) ).Get( []byte( key ) )
    dec := codec.NewBinaryCodecFrom( b )
    dec.Read( s )

    if s.TrainUID != "" {
      statistics.Incr( "schedule.200" )
      w.WriteHeader( 200 )
      json.NewEncoder( w ).Encode( s )
    } else {
      statistics.Incr( "schedule.404" )
      w.WriteHeader( 404 )
    }

    return nil
  }); err != nil {
    log.Println( "Get schedule", key, err )
    statistics.Incr( "schedule.500" )
    w.WriteHeader( 500 )
  }
}
