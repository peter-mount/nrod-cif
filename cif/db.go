package cif

import (
  "database/sql"
  "errors"
  "log"
)

// OpenDB opens a CIF database.
func (c *CIF) OpenDB( db *sql.DB ) error {
  if c.db != nil {
    return errors.New( "CIF Already attached to a Database" )
  }

  c.db = db

  return nil

}

// Close the database.
// If OpenDB() was used to open the db then that db is closed.
// If UseDB() was used this simply detaches the CIF from that DB. The DB is not closed()
func (c *CIF) Close() {

  // Only close if we own the DB, e.g. via OpenDB()
  if c.db != nil {
    c.db.Close()
  }

  // Detach
  c.db = nil
}

func (c *CIF) Execute( label, cmd string ) error {
  log.Println( label )
  _, err := c.db.Exec( cmd )
  return err
}

func (c *CIF) Cluster() error {
  log.Println( "Clustering station index" )

  _, err := c.db.Exec( "CLUSTER timetable.station USING station_tdt" )
  if err != nil {
    return err
  }

  log.Println( "Compacting station index" )
  _, err = c.db.Exec( "VACUUM FULL timetable.station" )
  return err
}

func (c *CIF) Cleanup( vacuum bool ) error {
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

  if vacuum || rc > 0 {
    log.Println( "Cleaning up freed space" )
    _, err = c.db.Exec( "VACUUM" )
    return err
  }

  return nil
}
