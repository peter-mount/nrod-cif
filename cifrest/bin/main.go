package main

import (
  "cifrest"
  "github.com/peter-mount/golib/kernel"
  "log"
)

func main() {
  err := kernel.Launch( &cifrest.CIFRest{} )
  if err != nil {
    log.Fatal( err )
  }
}
