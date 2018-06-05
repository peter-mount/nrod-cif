package cifrest

import (
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nrod-cif/cif"
  "log"
  "time"
)

// TiplocHandler implements a net/http handler that implements a simple Rest service to retrieve Tiploc records.
// The handler must have {id} set in the path for this to work, where id would represent the Tiploc code.
//
// For example:
//
// router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /tiploc/MSTONEE to return JSON representing that station.
func (c *CIFRest) ScheduleHandler( r *rest.Rest ) error {
  uid := r.Var( "uid" )
  date := r.Var( "date" )
  stp := r.Var( "stp" )

  startDate, err := time.Parse( "2006-01-02", date )
  if err != nil {
    r.Status( 500 )
    log.Printf( "500: schedule date %s %s %s = %s", uid, date, stp, err )
    return err
  }

  schedule, err := c.cif.GetSchedule( uid, startDate, stp )

  if err != nil {
    r.Status( 500 )
    log.Printf( "500: schedule %s %s %s = %s", uid, date, stp, err )
    return err
  }

  if schedule == nil {
    r.Status( 404 )
    return nil
  }

  resp := &cif.Response{
    TrainUID: uid,
    Date: date,
    STPIndicator: stp,
    Schedules: []*cif.Schedule{schedule},
    Self: r.Self( "/schedule/" + uid + "/" + date + "/" + stp ),
  }

  unknownTiplocs := resp.GetScheduleTiplocs( schedule )
  tiplocs, err := c.cif.GetTiplocs( unknownTiplocs )
  if err != nil {
    r.Status( 500 )
    log.Printf( "500: schedule tiplocs %s %s %s = %s", uid, date, stp, err )
    return err
  }

  resp.AddTiplocs( tiplocs )

  r.Status( 200 ).Value( resp )
  return nil
}
