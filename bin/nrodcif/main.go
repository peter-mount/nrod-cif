// CIF Rest server
package main


import (
//  "github.com/peter-mount/golib/rest"
  "bin"
  "cif"
)

func main() {
  bin.RunApplication( app )
}

func app( config *bin.Config ) ( func(), error ) {

  cif := &cif.CIF{}

  config.DbPath( &config.Database.Cif, "cif.db" )

  if err := cif.OpenDB( config.Database.Cif ); err != nil {
    return nil, err
  }

  return nil, nil
}
/*
  server.Handle( "/crs/{id}", db.CRSHandler ).Methods( "GET" )
  server.Handle( "/stanox/{id}", db.StanoxHandler ).Methods( "GET" )
  server.Handle( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )

  server.Handle( "/schedule/{uid}/{date}/{stp}", db.ScheduleHandler ).Methods( "GET" )
  server.Handle( "/schedule/{uid}", db.ScheduleUIDHandler ).Methods( "GET" )

  server.Handle( "/importCIF", db.ImportHandler ).Methods( "POST" )

  server.Start()
}
*/
