package cif

import (
  bolt "github.com/coreos/bbolt"
  "bytes"
  "encoding/json"
  "encoding/xml"
  "github.com/gorilla/mux"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/statistics"
  "log"
  "net/http"
)

// GetSchedulesByUID returns all Schedule's for a specific TrainUID.
// If no schedules exist for the required TrainUID then the returned slice is empty.
func (c *CIF) GetSchedulesByUID( tx *bolt.Tx, uid string ) []*Schedule {
  var ar []*Schedule

  b := tx.Bucket( []byte( "Schedule" ) ).Cursor()
  prefix := []byte( uid )

  for k, v := b.Seek( prefix ); k != nil && bytes.Compare( k[:len(prefix)], prefix ) == 0; k, v = b.Next() {
    s := &Schedule{}
    dec := codec.NewBinaryCodecFrom( v )
    dec.Read( s )
    if uid == s.TrainUID {
      ar = append( ar, s )
    }
  }

  return ar
}

// ScheduleUIDHandler implements a net/http handler that implements a simple Rest service to retrieve all schedules for a specific uid
// The handler must have {uid} set in the path for this to work.
//
// For example:
//
// router.HandleFunc( "/schedule/{uid}", db.ScheduleUIDHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct.
func (c *CIF) ScheduleUIDHandler( w http.ResponseWriter, r *http.Request ) {
  var params = mux.Vars( r )

  uid := params[ "uid" ]

  if err := c.db.View( func( tx *bolt.Tx ) error {

    ary := c.GetSchedulesByUID( tx, uid )
    if len( ary ) > 0 {
      statistics.Incr( "schedule.uid.200" )
      w.WriteHeader( 200 )

      if r.Header.Get( "Accept" ) == "text/xml" {
        // Wrap it in a schedules element
        var ret schedules = schedules{ Schedules: ary }
        xml.NewEncoder( w ).Encode( ret )
      } else {
        json.NewEncoder( w ).Encode( ary )
      }
    } else {
      statistics.Incr( "schedule.uid.404" )
      w.WriteHeader( 404 )
    }

    return nil
  }); err != nil {
    log.Println( "Get schedule", uid, err )
    statistics.Incr( "schedule.uid.500" )
    w.WriteHeader( 500 )
  }
}

// Wrapper used when writing the response as XML
type schedules struct {
    XMLName    xml.Name     `xml:"schedules"`
    Schedules  []*Schedule  `xml:"schedule"`
}
