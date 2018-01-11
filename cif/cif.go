// NR CIF file format

package cif

import (
  bolt "github.com/coreos/bbolt"
)

type CIF struct {
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
