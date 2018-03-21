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

  if config.NetworkRail.User.Username != "" && config.NetworkRail.User.Password != "" {
    cif.SetUpdater( config.NetworkRail.User.Username, config.NetworkRail.User.Password )
  }

  config.DbPath( &config.Database.Cif, "cif.db" )

  if err := cif.OpenDB( config.Database.Cif ); err != nil {
    return nil, err
  }

  cif.InitRest( config.Server.Ctx )

  return nil, nil
}
