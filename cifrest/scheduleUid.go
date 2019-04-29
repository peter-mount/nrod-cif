package cifrest

import (
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nrod-cif/cif"
  "log"
)

// ScheduleUIDHandler implements a REST endpoint which returns all timetable
// schedules for a specific schedule uid.
//
// Handler setup:
//
// router.HandleFunc( "/schedule/{uid}", a.ScheduleUIDHandler ).Methods( "GET" )
//
// where a is a pointer to an active CIF struct.
func (c *CIFRest) ScheduleUIDHandler(r *rest.Rest) error {
  uid := r.Var("uid")

  schedules, err := c.cif.GetSchedulesByUID(uid)

  if err != nil {
    r.Status(500)
    log.Printf("500: schedule %s = %s", uid, err)
    return err
  }

  if schedules == nil {
    r.Status(404)
    return nil
  }

  resp := &cif.Response{
    TrainUID:  uid,
    Schedules: schedules,
    Self:      r.Self("/schedule/" + uid),
  }

  for _, schedule := range schedules {
    unknownTiplocs := resp.GetScheduleTiplocs(schedule)
    tiplocs, err := c.cif.GetTiplocs(unknownTiplocs)
    if err != nil {
      r.Status(500)
      log.Printf("500: schedule tiplocs %s = %s", uid, err)
      return err
    }

    resp.AddTiplocs(tiplocs)
  }

  r.Status(200).Value(resp)
  return nil
}
