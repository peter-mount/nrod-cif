package cifimport

import (
  "log"
)

func (c *CIFImporter) cleanup( vacuum bool ) error {
  log.Println( "Removing historic schedules" )
  res, err := c.db.Exec( "DELETE FROM timetable.schedule WHERE enddate < NOW()::DATE" )
  if err != nil {
    return err
  }
  rc, err := res.RowsAffected()
  if err != nil {
    return err
  }
  if rc > 0 {
    log.Printf( "Removed %d schedules", rc )
  }

  if vacuum {
    log.Println( "Compacting tiplocs" )
    _, err = c.db.Exec( "VACUUM FULL timetable.tiploc" )
    if err != nil {
      return err
    }

    log.Println( "Compacting schedules" )
    _, err = c.db.Exec( "VACUUM FULL timetable.schedule" )
    if err != nil {
      return err
    }
    _, err = c.db.Exec( "VACUUM FULL timetable.schedule_json" )
    if err != nil {
      return err
    }

    log.Println( "Compacting station index" )
    _, err = c.db.Exec( "VACUUM FULL timetable.station" )
    if err != nil {
      return err
    }
  }

  if vacuum || rc > 0 {
    log.Println( "Cleaning up freed space" )
    _, err = c.db.Exec( "VACUUM" )
    return err
  }

  return nil
}
