// A GO library providing a database based on the Network Rail CIF Timetable feed.
package cif

import (
  "database/sql"
)

type CIF struct {
  // The DB
  db           *sql.DB
  // Last import HD record
  header       *HD
  // Current import HD record
  importhd     *HD
  // === Entries used during import only
  tx           *sql.Tx
  //
  curSchedule  *Schedule
  update        bool
  //
  //Updater     *Updater
  Timetable    CursorUpdate
}

type CursorUpdate interface {
  //Update( *CIF, *bolt.Bucket ) error
}

// String returns a human readable description of the latest CIF file imported into this database.
func (c *CIF) String() string {
  return c.header.String()
}

func (c *CIF) Update( f func( *sql.Tx ) error ) error {
  tx, err := c.db.Begin()
  if err != nil {
    return err
  }
  defer tx.Commit()

  err = f( tx )
  if err != nil {
    tx.Rollback()
    return err
  }

  return nil
}

func (c *CIF) UpdateTimetable() error {
  /*
  if c.Timetable != nil {
    return c.db.View( func( tx *bolt.Tx ) error {
      bucket := tx.Bucket( []byte("Schedule") )
      return c.Timetable.Update( c, bucket )
    })
  }
  */
  return nil
}
