package cif

import (
  "database/sql"
  "errors"
)

// OpenDB opens a CIF database.
func (c *CIF) OpenDB( db *sql.DB ) error {
  if c.db != nil {
    return errors.New( "CIF Already attached to a Database" )
  }

  c.db = db

  return nil

}
