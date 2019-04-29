package main

import (
  "github.com/peter-mount/golib/kernel"
  "github.com/peter-mount/nrod-cif/cifretrieve"
  "log"
)

func main() {
  err := kernel.Launch(&cifretrieve.CIFRetriever{})
  if err != nil {
    log.Fatal(err)
  }
}
