// A GO library providing a database based on the Network Rail CIF Timetable feed.
package cif

import (
  bolt "github.com/coreos/bbolt"
)

// Bitmasks for CIF.Mode used by CIF.Parse() & CIF.ParseFile() to determine what
// to import. If not set then everything is imported
const (
  // Import tiplocs only
  TIPLOC    = 1
  // Import schedules only
  SCHEDULE  = 1<<1
  // The default mode used if nothing is set
  ALL       = TIPLOC | SCHEDULE
)

type CIF struct {
  // The mode the parser should use when importing NR CIF files.
  // This is a bit mask of TIPLOC or SCHEDULE. If not set then ALL is used.
  Mode          int
  // Allow CIF.Close() to close the database.
  allowClose    bool
  // The DB
  db           *bolt.DB
  // Last import HD record
  header       *HD
  // Current import HD record
  importhd     *HD
  // === Entries used during import only
  tx           *bolt.Tx
  //
  curSchedule  *Schedule
  update        bool
  //
  tiploc       *bolt.Bucket
  crs          *bolt.Bucket
  stanox       *bolt.Bucket
  schedule     *bolt.Bucket
}

// String returns a human readable description of the latest CIF file imported into this database.
func (c *CIF) String() string {
  return c.header.String()
}
