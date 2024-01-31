package main

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nrod-cif/tools/cifimport"
	"log"
)

func main() {
	err := kernel.Launch(&kernel.MemUsage{}, &cifimport.CIFImporter{})
	if err != nil {
		log.Fatal(err)
	}
}
