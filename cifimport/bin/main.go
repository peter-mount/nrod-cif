package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nrod-cif/cifimport"
  "log"
)

func main() {
  err := kernel.Launch( &kernel.MemUsage{}, &cifimport.CIFImporter{} )
  if err != nil {
    log.Fatal( err )
  }
}
