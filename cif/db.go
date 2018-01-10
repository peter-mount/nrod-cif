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
    c.Header = &HD{}
    c.schedules = make( map[string][]*Schedule )

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

  return c, nil
}

// Ensures we have the appropriate buckets
func (c *CIF) initDB() error {
  tx, err := c.db.Begin(true)
  if err != nil {
    return err
  }
  defer tx.Rollback()

  var rebuildRequired bool

  for _, n := range []string { "Tiploc", "Crs", "Stanox" } {
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
    if err := c.Rebuild( tx ); err != nil {
      return err
    }
  }

  return tx.Commit()
}
