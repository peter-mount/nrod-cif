package cif

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
)

// ScheduleHandler implements a net/http handler that implements a simple Rest
// service to retrieve all schedules for a specific uid, date and STPIndicator
// The handler must have {uid} set in the path for this to work.
//
// For example:
//
// router.HandleFunc( "/schedule/{uid}/{date}/{stp}", db.ScheduleHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct.
func (c *CIF) ScheduleHandler( r *rest.Rest ) error {
  return c.db.View( func( tx *bolt.Tx ) error {

    uid := r.Var( "uid" )
    date := r.Var( "date" )
    stp := r.Var( "stp" )
    key := uid + date + stp

    s := &Schedule{}
    b := tx.Bucket( []byte( "Schedule" ) ).Get( []byte( key ) )
    dec := codec.NewBinaryCodecFrom( b )
    dec.Read( s )

    if s.ID.TrainUID != "" {
      statistics.Incr( "schedule.200" )
      result := NewResponse()
      result.Schedules = []*Schedule{s}
      c.ResolveScheduleTiplocs( tx, s, result)
      result.TiplocsSetSelf( r )
      s.SetSelf( r )
      result.SetSelf( r, "/schedule/" + uid + "/" + date + "/" + stp )
    } else {
      statistics.Incr( "schedule.404" )
      r.Status( 404 )
    }

    return nil
  })
}
