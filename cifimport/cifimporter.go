// CIF Importer
package cifimport

import (
  "cif"
  "database/sql"
  "flag"
  "fmt"
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/golib/kernel/db"
  "log"
  "os"
)

type CIFImporter struct {
  files   []string
  // The DB
  dbService    *db.DBService
  db           *sql.DB
  // Last import HD record
  header       *HD
  // Current import HD record
  importhd     *HD
  // === Entries used during import only
  tx           *sql.Tx
  //
  curSchedule  *cif.Schedule
  update        bool
  // Maintenance Mode
  maintenance  *bool
  forceExpire  *bool
  forceVacuum  *bool
}

func (a *CIFImporter) Name() string {
  return "CIFImporter"
}

func (a *CIFImporter) Init( k *kernel.Kernel ) error {
  a.maintenance = flag.Bool( "m", false, "Same as -expire -vacuum" )
  a.forceExpire = flag.Bool( "expire", false, "Remove expired entries" )
  a.forceVacuum = flag.Bool( "vacuum", false, "Vacuum & recluster the database" )

  dbservice, err := k.AddService( &db.DBService{} )
  if err != nil {
    return err
  }

  a.dbService = (dbservice).(*db.DBService)
  return nil
}

func (a *CIFImporter) PostInit() error {

  // Fail if we have no CIF files in the command line
  a.files = flag.Args()

  if *(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum) {
    if len( a.files ) > 0 {
      return fmt.Errorf( "CIF files not permitted in maintenance mode" )
    }
  } else if len( a.files ) == 0 {
    return fmt.Errorf( "CIF files required" )
  }

  return nil
}

func (a *CIFImporter) Start() error {
  a.db = a.dbService.GetDB()
  if a.db == nil {
    return fmt.Errorf( "No database" )
  }
  return nil
}

func (a *CIFImporter) Run() error {

  fileCount := 0

  // Normal mode, cleanup & import CIF files
  if !( *(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum) ) {
    // Do a cleanup first as it will remove expired entries freeing up some space
    err := a.cleanup()
    if err != nil {
      return err
    }

    for _, file := range a.files {

      log.Printf( "Parsing %s", file )

      f, err := os.Open( file )
      if err != nil {
        return err
      }
      defer f.Close()

      skip, err := a.importCIF( f )
      if err != nil {
        if skip {
          // Non fatal error so log it but don't kill the import
          log.Println( err )
        } else {
          return err
        }
      } else {
        fileCount ++;
      }
    }
  }

  // Do maintenance if in maintenance mode or we imported at least 1 CIF file
  if *(a.maintenance) || *(a.forceExpire) || fileCount > 0 {
    err := a.cleanup()
    if err != nil {
      return err
    }
  }

  if *(a.maintenance) || *(a.forceVacuum) || fileCount > 0 {
    err := a.vacuum()
    if err != nil {
      return err
    }

    err = a.cluster()
    if err != nil {
      return err
    }
  }

  if *(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum) {
    log.Println( "Maintenance complete" )
  } else {
    log.Println( "Import complete" )
  }

  return nil
}

func (c *CIFImporter) Update( f func( *sql.Tx ) error ) error {
  tx, err := c.db.Begin()
  if err != nil {
    return err
  }
  defer tx.Commit()

  err = f( tx )
  if err != nil {
    tx.Rollback()
    return err
  }

  return nil
}
