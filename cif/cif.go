// NR CIF file format

package cif

import (
  bolt "github.com/coreos/bbolt"
)

// Bitmasks for CIF.Mode used by CIF.Parse() & CIF.ParseFile() to determine what
// to import. If not set then everything is imported
const (
  TIPLOC    = 1     // Import tiplocs only
  SCHEDULE  = 1<<1  // Import schedules only
  // The default mode
  ALL       = TIPLOC | SCHEDULE
)

type CIF struct {
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

func (c *CIF) String() string {
  return c.header.String()
}
