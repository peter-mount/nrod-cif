// cifrest A standalone rest server for service CIF timetables
package cifrest

import (
  "fmt"
  cifdb "github.com/peter-mount/nrod-cif/cif/db"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/db"
  "github.com/peter-mount/golib/rest"
)

type CIFRest struct {
  cif           cifdb.CIF
  dbService    *db.DBService
  restService  *rest.Server
}

func (a *CIFRest) Name() string {
  return "CIFImporter"
}

func (a *CIFRest) Init( k *kernel.Kernel ) error {
  service, err := k.AddService( &db.DBService{} )
  if err != nil {
    return err
  }
  a.dbService = (service).(*db.DBService)

  service, err = k.AddService( &rest.Server{} )
  if err != nil {
    return err
  }
  a.restService = (service).(*rest.Server)

  return nil
}

func (a *CIFRest) PostInit() error {
  a.restService.Handle( "/crs/{id}", a.CRSHandler ).Methods( "GET" )
  a.restService.Handle( "/stanox/{id}", a.StanoxHandler ).Methods( "GET" )
  a.restService.Handle( "/tiploc/{id}", a.TiplocHandler ).Methods( "GET" )

  a.restService.Handle( "/schedule/{uid}/{date}/{stp}", a.ScheduleHandler ).Methods( "GET" )
  a.restService.Handle( "/schedule/{uid}", a.ScheduleUIDHandler ).Methods( "GET" )

  a.restService.Handle( "/timetable/{crs}/{date}", a.TimetableHandler ).Methods( "GET" )

  return nil
}

func (a *CIFRest) Start() error {
  d := a.dbService.GetDB()
  if d == nil {
    return fmt.Errorf( "No database" )
  }
  return a.cif.OpenDB( d )
}
