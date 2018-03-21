package cif

import (
  bolt "github.com/coreos/bbolt"
  "errors"
  "log"
  "time"
)

// OpenDB opens a CIF database.
func (c *CIF) OpenDB( dbFile string ) error {
  if c.db != nil {
    return errors.New( "CIF Already attached to a Database" )
  }

  if db, err := bolt.Open( dbFile, 0666, &bolt.Options{
    Timeout: 5 * time.Second,
    } ); err != nil {
      return err
  } else {
    c.db = db

    // Set the default mode for the parser
    if ( c.Mode & ALL ) == 0 {
      c.Mode = ALL
    }

    // Now ensure the DB is initialised with the required buckets
    if err := c.initDB(); err != nil {
      return err
    }

    if h, err := c.GetHD(); err != nil {
      return err
    } else {
      c.header = h

      if h.Id == "" {
        log.Println( "NOTICE: Database requires a full CIF import" )
        if c.Updater != nil {
          go c.Updater.Update()
        }
      } else {
        log.Println( "Database:", h )
      }
    }
  }

  return nil

}

// Close the database.
// If OpenDB() was used to open the db then that db is closed.
// If UseDB() was used this simply detaches the CIF from that DB. The DB is not closed()
func (c *CIF) Close() {

  // Only close if we own the DB, e.g. via OpenDB()
  if c.allowClose && c.db != nil {
    c.db.Close()
  }

  // Detach
  c.db = nil
}

// Ensures we have the appropriate buckets
func (c *CIF) initDB() error {

  buckets := []string { "Meta" }

  if (c.Mode & TIPLOC) == TIPLOC {
    buckets = append( buckets, "Tiploc", "Crs", "Stanox" )
  }

  if (c.Mode & SCHEDULE) == SCHEDULE {
    buckets = append( buckets, "Schedule" )
  }

  return c.db.Update( func( tx *bolt.Tx ) error {

    for _, n := range buckets {
      var nb []byte = []byte(n)
      if bucket := tx.Bucket( nb ); bucket == nil {
        log.Println( "Creating bucket", n )
        if _, err := tx.CreateBucket( nb ); err != nil {
          return err
        }
      }
    }

    return nil
  })
}

// Clear out a bucket
func (c *CIF) clearBucket( bucket *bolt.Bucket ) error {
  return bucket.ForEach( func( k, v []byte) error {
    return bucket.Delete( k )
  })
}

// Used in full imports, clears the relevant buckets
func (c *CIF) resetDB() error {

  if (c.Mode & TIPLOC) == TIPLOC {
    if err := c.clearBucket( c.tiploc ); err != nil {
      return err
    }

    if err := c.clearBucket( c.crs ); err != nil {
      return err
    }

    if err := c.clearBucket( c.stanox ); err != nil {
      return err
    }
  }

  if (c.Mode & SCHEDULE) == SCHEDULE {
    if err := c.clearBucket( c.schedule ); err != nil {
      return err
    }
  }

  return nil
}
