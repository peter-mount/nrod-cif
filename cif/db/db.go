package db

import (
  "database/sql"
  "errors"
)

type CIF struct {
  // The DB
  db           *sql.DB
}

// OpenDB opens a CIF database.
func (c *CIF) OpenDB( d *sql.DB ) error {
  if c.db != nil {
    return errors.New( "CIF Already attached to a Database" )
  }

  c.db = d

  return nil
}
