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

  cif.InitRest( config.Server.Ctx )

  return nil, nil
}
