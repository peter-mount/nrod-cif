// CIF Rest server
package main

import (
  bolt "github.com/coreos/bbolt"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
)

func (c *CifDB) OpenDB( dbFile string ) error {

  if boltdb, err := bolt.Open( dbFile, 0666, &bolt.Options{
    Timeout: 5 * time.Second,
    } ); err != nil {
      return err
  } else {
    db.db = boltdb
    db.cif.Init( db.db )
  }

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )
    db.db.Close()
    log.Println( "Database closed" )
    os.Exit( 0 )
  }()

  // Now ensure the DB is initialised with the required buckets
  if err := initDB(); err != nil {
    return err
  }

  return nil
}

// Ensures we have the appropriate buckets
func initDB() error {
  tx, err := db.db.Begin(true)
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
    if err := db.cif.Rebuild( tx ); err != nil {
      return err
    }
  }

  return tx.Commit()
}
