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
}

func (a *CIFImporter) Name() string {
  return "CIFImporter"
}

func (a *CIFImporter) Init( k *kernel.Kernel ) error {
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
  if len( a.files ) == 0 {
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

  // Do a cleanup first
  err := a.cleanup( false )
  if err != nil {
    return err
  }

  fileCount := 0

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

  if fileCount > 0 {
    err = a.cleanup( true )
    if err != nil {
      return err
    }

    err = a.cluster()
    if err != nil {
      return err
    }
  }

  log.Println( "Import complete" )
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
