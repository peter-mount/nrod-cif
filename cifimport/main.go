// CIF Importer
package main

import (
  "github.com/peter-mount/golib/kernel"
  "log"
)

func main() {
  err := kernel.Launch( &CIFImporter{} )
  if err != nil {
    log.Fatal( err )
  }
}
