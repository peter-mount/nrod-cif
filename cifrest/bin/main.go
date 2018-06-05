package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nrod-cif/cifrest"
  "log"
)

func main() {
  err := kernel.Launch( &cifrest.CIFRest{} )
  if err != nil {
    log.Fatal( err )
  }
}
