// Handles the persistance of the imported CIF file
package cif

import (
  "io"
  "encoding/gob"
)

// Copy of CIF but with fields public - required for gob to work
type dbFile struct {
  Header   *HD
  Tiploc    map[string]*Tiploc
  // Map of Schedules
  Schedules map[string][]*Schedule
}

func (c *CIF) Write( w io.Writer ) error {
  var db dbFile
  db.Header = c.Header
  db.Tiploc = c.tiploc
  db.Schedules = c.schedules

  enc := gob.NewEncoder( w )

  return enc.Encode( db )
}

func (c *CIF) Read( r io.Reader ) error {
  var db dbFile

  dec := gob.NewDecoder( r )
  if err := dec.Decode( &db ); err != nil {
    return err
  }

  c.Header = db.Header
  c.tiploc = db.Tiploc
  c.schedules = db.Schedules
  c.cleanup()
  return nil
}
