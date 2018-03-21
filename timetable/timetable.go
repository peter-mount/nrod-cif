package timetable

import (
  bolt "github.com/coreos/bbolt"
  "errors"
  "log"
  "time"
)

// A searchable timetable based on a CIF
type Timetable struct {
  // The DB for the Timetable
  db           *bolt.DB
}

// OpenDB opens a CIF database.
func (c *Timetable) OpenDB( dbFile string ) error {
  if c.db != nil {
    return errors.New( "CIF Already attached to a Database" )
  }

  if db, err := bolt.Open( dbFile, 0666, &bolt.Options{
    Timeout: 5 * time.Second,
    } ); err != nil {
      return err
  } else {
    c.db = db

    // Now ensure the DB is initialised with the required buckets
    if err := c.initDB(); err != nil {
      return err
    }
  }

  return nil

}

// Close the database.
// If OpenDB() was used to open the db then that db is closed.
// If UseDB() was used this simply detaches the CIF from that DB. The DB is not closed()
func (c *Timetable) Close() {

  // Only close if we own the DB, e.g. via OpenDB()
  if c.db != nil {
    c.db.Close()
  }

  // Detach
  c.db = nil
}

// Ensures we have the appropriate buckets
func (c *Timetable) initDB() error {

  buckets := []string { "Timetable" }

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
func (c *Timetable) clearBucket( bucket *bolt.Bucket ) error {
  return bucket.ForEach( func( k, v []byte) error {
    return bucket.Delete( k )
  })
}
