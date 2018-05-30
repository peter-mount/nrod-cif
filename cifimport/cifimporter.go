// CIF Importer
package main

import (
  "cifservice"
  "flag"
  "fmt"
  "github.com/peter-mount/golib/kernel"
  "log"
  "os"
)

type CIFImporter struct {
  cif      *cifservice.CIFService
  files   []string
}

func (a *CIFImporter) Name() string {
  return "CIFImporter"
}

func (a *CIFImporter) Init( k *kernel.Kernel ) error {
  s, err := k.AddService( &cifservice.CIFService{} )
  if err != nil {
    return err
  }

  a.cif = s.(*cifservice.CIFService)

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

func (a *CIFImporter) Run() error {

  // Do a cleanup first
  err := a.cif.Cif.Cleanup( false )
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

    skip, err := a.cif.Cif.ImportCIF( f )
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
    err = a.cif.Cif.Cleanup( true )
    if err != nil {
      return err
    }

    err = a.cif.Cif.Cluster()
    if err != nil {
      return err
    }
  }

  log.Println( "Import complete" )
  return nil
}
