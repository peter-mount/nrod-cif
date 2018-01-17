package cif

import (
  bolt "github.com/coreos/bbolt"
  "bytes"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
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
func (c *CIF) ScheduleUIDHandler( r *rest.Rest ) error {
  return c.db.View( func( tx *bolt.Tx ) error {

    uid := r.Var( "uid" )

    result := NewResponse()
    r.Value( result )

    result.Schedules = c.GetSchedulesByUID( tx, uid )
    if len( result.Schedules ) > 0 {
      statistics.Incr( "schedule.uid.200" )
      r.Status( 200 )
      result.Status = 200
      result.Self = r.Self( "/schedule/" + uid )
      for _, s := range result.Schedules {
        c.ResolveScheduleTiplocs( tx, s, result )
        s.SetSelf( r )
        result.TiplocsSetSelf( r )
      }
    } else {
      statistics.Incr( "schedule.uid.404" )
      r.Status( 404 )
      result.Status = 404
    }

    return nil
  })
}
