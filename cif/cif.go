// A GO library providing a database based on the Network Rail CIF Timetable feed.
package cif

import (
  "database/sql"
)

type CIF struct {
  // The DB
  db           *sql.DB
}

const (
  DateTime        = "2006-01-02 15:04:05"
  Date            = "2006-01-02"
  HumanDateTime   = "2006 Jan 02 15:04:05"
  HumanDate       = "2006 Jan 02"
  Time            = "15:04:05"
)

// Scanable is an interface with a function called Scan.
// e.g. Row or Rows in database/sql
type Scannable interface {
  Scan(... interface{}) error
}
