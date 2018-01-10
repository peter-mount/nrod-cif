// CIF Import utility

package main

import (
  "cif"
  "flag"
  "log"
)

func main() {
  log.Println( "cifimport v0.1" )

  dbFile := flag.String( "f", "", "The CIF file to manage" )

  flag.Parse()

  if dbFile == nil || *dbFile == "" {
    log.Println( "Source CIF file -f required" )
  }

  // Our copy of the CIF file
  db := &cif.CIF{}
  db.Init()

  // todo load here

  // Now parse each file
  for _, arg := range flag.Args() {
    log.Println( "Importing", arg )

    err := db.Parse( arg )
    if err != nil {
      log.Fatal( "Failed to parse", arg, err )
    }

  }

  log.Println( "Import complete" )
  log.Println( db )

  if s := db.GetSchedules( "Y74216" ); s != nil {
    for _, sched := range s {
      log.Println( sched.FullString() )
    }
  }

  // todo write file
}
