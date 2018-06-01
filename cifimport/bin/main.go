package main

import (
  "cifimport"
  "github.com/peter-mount/golib/kernel"
  "log"
)

func main() {
  err := kernel.Launch( &cifimport.CIFImporter{} )
  if err != nil {
    log.Fatal( err )
  }
}
