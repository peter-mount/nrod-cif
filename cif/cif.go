// NR CIF file format

package cif

import (
  bolt "github.com/coreos/bbolt"
  "github.com/pkg/errors"
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

func (c *CIF) get( b *bolt.Bucket, k string, i interface{} ) error {
  bar := b.Get( []byte(k) )
  if bar != nil {
    if err := getInterface( bar, i ); err != nil {
      return errors.WithStack( err )
    }
    return nil
  }
  return errors.New( k + " Not found")
}

func (c *CIF) put( b *bolt.Bucket, k string, i interface{} ) error {
  if bar, err := getBytes( i ); err != nil {
    return err
  } else {
    if err := b.Put( []byte(k), bar ); err != nil {
      return errors.WithStack( err )
    }
  }
  return nil
}

func (c *CIF) resetDB() error {
  if err := c.tiploc.ForEach( func( k, v []byte) error {
    return c.tiploc.Delete( k )
  }); err != nil {
    return err
  }

  return c.schedule.ForEach( func( k, v []byte) error {
    return c.schedule.Delete( k )
  })
}

func (c *CIF) String() string {
  return c.header.String()
}

func (c *CIF) Rebuild( tx *bolt.Tx ) error {

  c.tiploc = tx.Bucket( []byte("Tiploc") )
  c.crs = tx.Bucket( []byte("Crs") )
  c.stanox = tx.Bucket( []byte("Stanox") )
  c.schedule = tx.Bucket( []byte("Schedule") )

  if err := c.cleanupStanox(); err != nil {
    return err
  }

  if err := c.cleanupCRS(); err != nil {
    return err
  }

  /*
  if err := c.cleanupSchedules(); err != nil {
    return err
  }
  */

  return nil
}
