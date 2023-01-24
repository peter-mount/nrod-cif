package main

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nrod-cif/cifimport"
	"os"
)

func main() {
	err := kernel.Launch(&kernel.MemUsage{}, &cifimport.CIFImporter{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
