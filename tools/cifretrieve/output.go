package cifretrieve

import (
	//  "fmt"
	//  "io"
	//  "io/ioutil"
	"log"
	//  "net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

func (a *CIFRetriever) writeoutput() error {
	log.Println("Writing output")

	// Sort by file date
	sort.SliceStable(a.files, func(i, j int) bool {
		return a.files[i].date.Before(a.files[j].date)
	})

	// Remove all before the last full export
	fullIndex := -1
	for i, cif := range a.files {
		if cif.full {
			fullIndex = i
		}
	}
	if fullIndex >= 0 {
		a.files = a.files[fullIndex:]

		// Remove any updates on the same date as the full import
		// index 0 will be the full, 1 the next update
		// What this does is ensure we don't run an update for the same day as the
		// full import - as it's just a waste of time as the full will contain the
		// update
		if len(a.files) > 1 {
			d0 := a.files[0].date.Truncate(24 * time.Hour)
			d1 := a.files[1].date.Truncate(24 * time.Hour)
			if d0.Equal(d1) {
				if len(a.files) == 2 {
					a.files = a.files[:1]
				} else {
					a.files = append(a.files[:1], a.files[2:]...)
				}
			}
		}
	}

	// Now write the files to the output file & log
	err := os.MkdirAll(filepath.Dir(*a.output), 0755)
	if err != nil {
		return err
	}

	wfile, err := os.OpenFile(*a.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer wfile.Close()

	for _, cif := range a.files {
		log.Println(cif.path)
		wfile.WriteString(cif.path + "\n")
	}

	return nil
}
