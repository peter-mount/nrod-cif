// NR CIF file format

package cif

import (
  bolt "github.com/coreos/bbolt"
  "errors"
  "fmt"
)

type CIF struct {
  db           *bolt.DB
  // === Entries used during import only
  tx           *bolt.Tx
  Header       *HD
  curSchedule  *Schedule
  tiploc       *bolt.Bucket
  crs          *bolt.Bucket
  stanox       *bolt.Bucket
  schedule     *bolt.Bucket
}

func (c *CIF) get( b *bolt.Bucket, k string, i interface{} ) error {
  bar := b.Get( []byte(k) )
  if bar != nil {
    return getInterface( bar, i )
  }
  return errors.New( k + " Not found")
}

func (c *CIF) put( b *bolt.Bucket, k string, i interface{} ) error {
  if bar, err := getBytes( i ); err != nil {
    return err
  } else {
    return b.Put( []byte(k), bar )
  }
}

func (c *CIF) resetDB() error {
  if err := c.tiploc.ForEach( func( k, v []byte) error {
    return c.tiploc.Delete( k )
  }); err != nil {
    return err
  }

  return c.schedule.ForEach( func( k, v []byte) error {
    return c.tiploc.Delete( k )
  })
}

func (c *CIF) String() string {
  return fmt.Sprintf(
    "CIF %s Extracted %v Date Range %v - %v Update %s",
    c.Header.FileMainframeIdentity,
    c.Header.DateOfExtract.Format( HumanDateTime ),
    c.Header.UserStartDate.Format( HumanDate ),
    c.Header.UserEndDate.Format( HumanDate ),
    c.Header.Update )
}

func (c *CIF) Rebuild( tx *bolt.Tx ) error {

  c.tiploc = tx.Bucket( []byte("Tiploc") )
  c.crs = tx.Bucket( []byte("Crs") )
  c.stanox = tx.Bucket( []byte("Stanox") )

  if err := c.cleanupStanox(); err != nil {
    return err
  }

  if err := c.cleanupCRS(); err != nil {
    return err
  }

  if err := c.cleanupSchedules(); err != nil {
    return err
  }

  return nil
}
