// CIF Rest server
package main

import (
  bolt "github.com/coreos/bbolt"
  "cif"
  "flag"
  "log"
)

type CifDB struct {
  db     *bolt.DB
  server  Server
  cif     cif.CIF
}

var db CifDB

func main() {
  log.Println( "cifserver v0.1" )

  writeSecret := flag.String( "s", "", "The write secret" )
  dbFile := flag.String( "d", "/database.db", "The database file" )
  port := flag.Int( "p", 8080, "Port to use" )
  flag.Parse()

  db.server.Port = *port
  log.Println( "secret", *writeSecret )
  log.Println( "dbFile", *dbFile )

  db.server.Init()

  if err := db.OpenDB( *dbFile ); err != nil {
    log.Fatal( err )
  }
  defer db.db.Close()

  if err := initDB(); err != nil {
    log.Fatal( err )
  }

  db.server.Router.HandleFunc( "/crs/{id}", db.cif.CRSHandler ).Methods( "GET" )
  db.server.Router.HandleFunc( "/stanox/{id}", db.cif.StanoxHandler ).Methods( "GET" )
  db.server.Router.HandleFunc( "/tiploc/{id}", db.cif.TiplocHandler ).Methods( "GET" )

  db.server.Router.HandleFunc( "/importCIF", db.cif.ImportHandler ).Methods( "POST" )

  db.server.Start()
}
