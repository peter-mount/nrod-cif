// CIF Rest server
package main


import (
//  "github.com/peter-mount/golib/rest"
  "bin"
  "cif"
  "log"
  "timetable"
)

func main() {
  bin.RunApplication( app )
}

func app( config *bin.Config ) ( func(), error ) {

  cif := &cif.CIF{}

  if config.NetworkRail.User.Username != "" && config.NetworkRail.User.Password != "" {
    cif.SetUpdater( config.NetworkRail.User.Username, config.NetworkRail.User.Password )
  }

  // Open timetable database
  config.DbPath( &config.Database.CifTimetable, "ciftt.db" )
  cifTimetable := &timetable.Timetable{}
  if err := cifTimetable.OpenDB( config.Database.CifTimetable ); err != nil {
    return nil, err
  }

  // Open cif database
  config.DbPath( &config.Database.Cif, "cif.db" )
  cif.Timetable = cifTimetable
  if err := cif.OpenDB( config.Database.Cif ); err != nil {
    cifTimetable.Close()
    return nil, err
  }

  cif.InitRest( config.Server.Ctx )

  return func() {
    log.Println( "Closing timetable" )
    cifTimetable.Close()

    log.Println( "Closing cif" )
    cif.Close()
  }, nil
}
