// CIF Rest server
package main

import (
  "github.com/peter-mount/golib/statistics"
  "cif"
  "flag"
  "log"
)

func main() {
  log.Println( "cifserver v0.1" )

  writeSecret := flag.String( "s", "", "The write secret" )
  dbFile := flag.String( "d", "/database.db", "The database file" )
  port := flag.Int( "p", 8080, "Port to use" )
  flag.Parse()

  stats := statistics.Statistics{ Log: true }
  stats.Configure()

  log.Println( "secret", *writeSecret )
  log.Println( "dbFile", *dbFile )

  db, err := cif.OpenCIF( *dbFile )
  if err != nil {
    log.Fatal( err )
  }

  var server Server = Server{ Port: *port }
  server.Init()

  server.Router.HandleFunc( "/crs/{id}", db.CRSHandler ).Methods( "GET" )
  server.Router.HandleFunc( "/stanox/{id}", db.StanoxHandler ).Methods( "GET" )
  server.Router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )

  server.Router.HandleFunc( "/importCIF", db.ImportHandler ).Methods( "POST" )

  server.Start()
}
