package main

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2"
	"github.com/peter-mount/nrod-cif/cifretrieve"
	"os"
)

func main() {
	err := kernel.Launch(&cifretrieve.CIFRetriever{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
