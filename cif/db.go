// CIF Rest server
package cif

import (
  bolt "github.com/coreos/bbolt"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
)

func OpenCIF( dbFile string ) ( *CIF, error ) {

  var c *CIF = &CIF{}

  if boltdb, err := bolt.Open( dbFile, 0666, &bolt.Options{
    Timeout: 5 * time.Second,
    } ); err != nil {
      return nil, err
  } else {
    c.db = boltdb
  }

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )
    c.db.Close()
    log.Println( "Database closed" )
    os.Exit( 0 )
  }()

  // Now ensure the DB is initialised with the required buckets
  if err := c.initDB(); err != nil {
    return nil, err
  }

  if h, err := c.GetHD(); err != nil {
    return nil, err
  } else {
    c.header = h

    if h.Id == "" {
      log.Println( "NOTICE: Database requires a full CIF import" )
    } else {
      log.Println( "Database:", h )
    }
  }

  return c, nil
}

// Ensures we have the appropriate buckets
func (c *CIF) initDB() error {
  return c.db.Update( func( tx *bolt.Tx ) error {
    var rebuildRequired bool

    for _, n := range []string { "Meta", "Tiploc", "Crs", "Stanox", "Schedule" } {
      var nb []byte = []byte(n)
      if bucket := tx.Bucket( nb ); bucket == nil {
        log.Println( "Creating bucket", n )
        if _, err := tx.CreateBucket( nb ); err != nil {
          return err
        }
        rebuildRequired = true
      }
    }

    if rebuildRequired {
      return c.Rebuild( tx )
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

// Store a struct into a bucket
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

// Retrieve a struct from a bucket
// If the entry does not exist then this returns nil
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
  if err := c.clearBucket( c.tiploc ); err != nil {
    return err
  }

  if err := c.clearBucket( c.crs ); err != nil {
    return err
  }

  if err := c.clearBucket( c.stanox ); err != nil {
    return err
  }

  return c.clearBucket( c.schedule )
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
