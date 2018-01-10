// CIF Import utility

package main

import (
  "cif"
  "flag"
  "log"
  "os"
)

func main() {
  log.Println( "cifimport v0.1" )

  srcFile := flag.String( "i", "", "The CIF file to manage" )
  destFile := flag.String( "o", "", "The CIF file to write" )

  flag.Parse()

  if destFile == nil || *destFile == "" {
    log.Fatal( "Output DB file -o required" )
  }

  // Our copy of the CIF file
  db := &cif.CIF{}
  db.Init()

  // Load an existing file
  if srcFile != nil && *srcFile != "" {
    log.Println( "Loading", *srcFile )
    if in, err := os.Open( *srcFile ); err != nil {
      log.Fatal( err )
    } else {
      defer in.Close()

      if err = db.Read( in ); err != nil {
        log.Fatal( err )
      }

      log.Println( db )
    }
  }

  // Now parse each file
  for _, arg := range flag.Args() {
    log.Println( "Importing", arg )

    if err := db.Parse( arg ); err != nil {
      log.Fatal( "Failed to parse", arg, err )
    }

  }

  log.Println( "Import complete" )
  log.Println( db )

  log.Println( "Writing", *destFile )
  if out, err := os.Create( *destFile ); err != nil {
    log.Fatal( err )
  } else {
    defer out.Close()

    if err = db.Write( out ); err != nil {
      log.Fatal( err )
    }
  }

  log.Println( "Written" )
/*
  if s := db.GetSchedules( "Y74216" ); s != nil {
    for _, sched := range s {
      log.Println( sched.FullString() )
    }
  }
*/

}
