package cifrest

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

// TiplocHandler implements a net/http handler that implements a simple Rest service to retrieve Tiploc records.
// The handler must have {id} set in the path for this to work, where id would represent the Tiploc code.
//
// For example:
//
// router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /tiploc/MSTONEE to return JSON representing that station.
func (c *CIFRest) ScheduleUIDHandler( r *rest.Rest ) error {
  uid := r.Var( "uid" )

  schedules, err := c.cif.GetSchedulesByUID( uid )

  if err != nil {
    r.Status( 500 )
    log.Printf( "500: schedule %s = %s", uid, err )
    return err
  }

  if schedules == nil {
    r.Status( 404 )
    return nil
  }

  resp := &Response{
    TrainUID: uid,
    Schedules: schedules,
    Self: r.Self( "/schedule/" + uid ),
  }

  for _, schedule := range schedules {
    unknownTiplocs := resp.GetScheduleTiplocs( schedule )
    tiplocs, err := c.cif.GetTiplocs( unknownTiplocs )
    if err != nil {
      r.Status( 500 )
      log.Printf( "500: schedule tiplocs %s = %s", uid, err )
      return err
    }

    resp.AddTiplocs( tiplocs )
  }

  r.Status( 200 ).Value( resp )
  return nil
}
