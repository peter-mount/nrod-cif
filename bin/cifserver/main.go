// CIF Rest server
package main

import (
  "github.com/peter-mount/golib/statistics"
  "cif"
  "flag"
  "os"
  "os/signal"
  "syscall"
  "log"
)

func main() {
  log.Println( "cifserver v0.1" )

  // TODO use this to protect /importCIF endpoint
  //writeSecret := flag.String( "s", "", "The write secret" )

  dbFile := flag.String( "d", "/database.db", "The database file" )

  // Port for the webserver
  port := flag.Int( "p", 8080, "Port to use" )

  flag.Parse()

  stats := statistics.Statistics{ Log: true }
  stats.Configure()

  db := cif.CIF{}

  if err := db.OpenDB( *dbFile ); err != nil {
    log.Fatal( err )
  }

  // Listen to signals & close the db before exiting
  // SIGINT for ^C, SIGTERM for docker stopping the container
  sigs := make( chan os.Signal, 1 )
  signal.Notify( sigs, syscall.SIGINT, syscall.SIGTERM )
  go func() {
    sig := <-sigs
    log.Println( "Signal", sig )

    db.Close()
    log.Println( "Database closed" )

    os.Exit( 0 )
  }()

  var server Server = Server{ Port: *port }
  server.Init()

  server.Router.HandleFunc( "/crs/{id}", db.CRSHandler ).Methods( "GET" )
  server.Router.HandleFunc( "/stanox/{id}", db.StanoxHandler ).Methods( "GET" )
  server.Router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )

  server.Router.HandleFunc( "/schedule/{uid}/{date}/{stp}", db.ScheduleHandler ).Methods( "GET" )
  server.Router.HandleFunc( "/schedule/{uid}", db.ScheduleUIDHandler ).Methods( "GET" )

  server.Router.HandleFunc( "/importCIF", db.ImportHandler ).Methods( "POST" )

  server.Start()
}
