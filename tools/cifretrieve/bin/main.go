package main

import (
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nrod-cif/tools/cifretrieve"
	"log"
)

func main() {
	err := kernel.Launch(&cifretrieve.CIFRetriever{})
	if err != nil {
		log.Fatal(err)
	}
}
