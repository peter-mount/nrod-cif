package cifimport

import (
  "github.com/peter-mount/nrod-cif/cif"
  "encoding/json"
  "log"
)

func (c *CIFImporter) addSchedule() error {
  // Do nothing if we have no schedule to persist
  if c.curSchedule == nil {
    return nil
  }

  // get schedule & reset
  s := c.curSchedule
  c.curSchedule = nil

  // Link it to this CIF file & persist
  s.DateOfExtract = c.importhd.DateOfExtract

  sj, err := json.Marshal( s )
  if err != nil {
    return err
  }

  _, err = c.tx.Exec( "SELECT timetable.addschedule( $1 )", sj )
  if err != nil {
    log.Printf( "Entry that failed:\n%s", string(sj) )
  }
  return err
}

func (c *CIFImporter) deleteSchedule( s *cif.Schedule ) error {
  _, err := c.tx.Exec(
    "DELETE FROM timetable.schedule WHERE uid = $1 AND stp = $2 AND startdate = $3",
    s.ID.TrainUID,
    s.ID.STPIndicator,
    s.Runs.RunsFrom.Format( Date ),
  )
  return err
}
