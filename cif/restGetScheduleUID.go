package cif

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
)

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
    result.Schedules = c.GetSchedulesByUID( tx, uid )
    if len( result.Schedules ) > 0 {
      statistics.Incr( "schedule.uid.200" )
      for _, s := range result.Schedules {
        c.ResolveScheduleTiplocs( tx, s, result )
        s.SetSelf( r )
        result.TiplocsSetSelf( r )
      }
      result.SetSelf( r, "/schedule/" + uid )
    } else {
      statistics.Incr( "schedule.uid.404" )
      r.Status( 404 )
    }

    return nil
  })
}
