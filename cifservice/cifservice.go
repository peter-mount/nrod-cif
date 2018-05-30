// Kernel service to wrap the CIF database
package cifservice

import (
  "cif"
  "fmt"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/db"
)

type CIFService struct {
  Cif     cif.CIF
  db     *db.DBService
}

func (a *CIFService) Name() string {
  return "CIF"
}

func (a *CIFService) Init( k *kernel.Kernel ) error {
  dbservice, err := k.AddService( &db.DBService{} )
  if err != nil {
    return err
  }

  a.db = (dbservice).(*db.DBService)
  return nil
}

func (a *CIFService) Start() error {
  d := a.db.GetDB()
  if d == nil {
    return fmt.Errorf( "No database" )
  }
  return a.Cif.OpenDB( d )
}
