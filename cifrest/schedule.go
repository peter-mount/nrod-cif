package cifrest

import (
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nrod-cif/cif"
  "log"
  "time"
)

// ScheduleHandler implements a REST endpoint which returns the details of a
// single specific schedule within the timetable.
//
// Handler setup:
//
// router.HandleFunc( "/schedule/{uid}/{date}/{stp}", a.ScheduleHandler ).Methods( "GET" )
//
// where a is a pointer to an active CIF struct.
func (c *CIFRest) ScheduleHandler(r *rest.Rest) error {
  uid := r.Var("uid")
  date := r.Var("date")
  stp := r.Var("stp")

  startDate, err := time.Parse("2006-01-02", date)
  if err != nil {
    r.Status(500)
    log.Printf("500: schedule date %s %s %s = %s", uid, date, stp, err)
    return err
  }

  schedule, err := c.cif.GetSchedule(uid, startDate, stp)

  if err != nil {
    r.Status(500)
    log.Printf("500: schedule %s %s %s = %s", uid, date, stp, err)
    return err
  }

  if schedule == nil {
    r.Status(404)
    return nil
  }

  resp := &cif.Response{
    TrainUID:     uid,
    Date:         date,
    STPIndicator: stp,
    Schedules:    []*cif.Schedule{schedule},
    Self:         r.Self("/schedule/" + uid + "/" + date + "/" + stp),
  }

  unknownTiplocs := resp.GetScheduleTiplocs(schedule)
  tiplocs, err := c.cif.GetTiplocs(unknownTiplocs)
  if err != nil {
    r.Status(500)
    log.Printf("500: schedule tiplocs %s %s %s = %s", uid, date, stp, err)
    return err
  }

  resp.AddTiplocs(tiplocs)

  r.Status(200).Value(resp)
  return nil
}
