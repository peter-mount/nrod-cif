// CIF Rest server
package main


import (
//  "github.com/peter-mount/golib/rest"
  "bin"
  "cif"
  "github.com/hashicorp/consul/api"
  "log"
  "timetable"
)

func main() {
  bin.RunApplication( app )
}

func app( config *bin.Config ) ( func(), error ) {

  if *consul {
    client, err := api.NewClient(api.DefaultConfig())
    if err != nil {
      log.Fatal( err )
    }
    if err := config.NetworkRail.ReadConsul( client ); err != nil {
      log.Fatal( err )
    }

    log.Println( "ServiceRegister" )
    agent := client.Agent()
    service := &api.AgentServiceRegistration{
      ID: "test",
      Name: "nrod-cif",
      Port: 80,
      //Address: "",
    }
    if err := agent.ServiceRegister( service ); err != nil {
      log.Fatal( "Failed to create service")
    }
    log.Println( "ServiceRegistered" )

    f, err := app1( config )
    if err != nil {
      agent.ServiceDeregister( "test" )
      return nil, err
    }

    return func() {
      log.Println( "Deregistering service" )
      agent.ServiceDeregister( "test" )
      f()
      }, nil
  }

  // no consul
  f, err := app1( config )
  return f, err
}

func app1( config *bin.Config ) ( func(), error ) {
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
