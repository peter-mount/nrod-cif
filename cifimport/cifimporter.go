package cifimport

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/peter-mount/go-kernel/db"
	"github.com/peter-mount/nrod-cif/cif"
	_ "gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

type CIFImporter struct {
	files     []string
	dbService *db.DBService `kernel:"inject"`
	db        *sql.DB
	//sql       *sqlutils.SchemaImport
	header      *HD     // Last import HD record
	importhd    *HD     // Current import HD record
	tx          *sql.Tx // === Entries used during import only
	curSchedule *cif.Schedule
	update      bool
	// Maintenance Mode
	maintenance *bool   `kernel:"flag,m,Same as -expire -vacuum"`
	forceExpire *bool   `kernel:"flag,expire,Remove expired entries"`
	forceVacuum *bool   `kernel:"flag,vacuum,Vacuum and re-cluster database"`
	fileSource  *string `kernel:"flag,files,files containing cif files to import"`
}

func (a *CIFImporter) Name() string {
	return "CIFImporter"
}

/*
func (a *CIFImporter) Init(k *kernel.Kernel) error {
  a.maintenance = flag.Bool("m", false, "Same as -expire -vacuum")
  a.forceExpire = flag.Bool("expire", false, "Remove expired entries")
  a.forceVacuum = flag.Bool("vacuum", false, "Vacuum & recluster the database")
  a.fileSource = flag.String("files", "", "File containing cif files to import")

  dbservice, err := k.AddService(&db.DBService{})
  if err != nil {
    return err
  }
  a.dbService = (dbservice).(*db.DBService)

  sqlservice, err := k.AddService(sqlutils.NewSchemaImport("timetable", AssetString, AssetNames))
  if err != nil {
    return err
  }
  a.sql = (sqlservice).(*sqlutils.SchemaImport)

  return nil
}
*/

func (a *CIFImporter) PostInit() error {

	files, err := a.addCIFFilesForImport()
	if err != nil {
		return err
	}
	a.files = files

	if *(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum) {
		if len(a.files) > 0 {
			return fmt.Errorf("CIF files not permitted in maintenance mode")
		}
	} else if len(a.files) == 0 {
		return fmt.Errorf("CIF files required")
	}

	return nil
}

func (a *CIFImporter) addCIFFilesForImport() ([]string, error) {
	files := flag.Args()

	if *a.fileSource != "" {
		file, err := os.Open(*a.fileSource)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			files = append(a.files, scanner.Text())
		}
		err = scanner.Err()
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (a *CIFImporter) Start() error {
	a.db = a.dbService.GetDB()
	if a.db == nil {
		return fmt.Errorf("No database")
	}
	return nil
}

func (a *CIFImporter) Run() error {

	fileCount := 0

	// Normal mode, cleanup & import CIF files
	if !(*(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum)) {
		// Do a cleanup first as it will remove expired entries freeing up some space
		err := a.cleanup()
		if err != nil {
			return err
		}

		for _, file := range a.files {

			log.Printf("Parsing %s", file)

			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			// gzip or plain
			var header [2]byte
			c, err := io.ReadFull(f, header[:])
			if err != nil {
				return err
			}
			if c < 2 {
				return fmt.Errorf("")
			}
			_, err = f.Seek(0, io.SeekStart)
			if err != nil {
				return err
			}
			reader := io.Reader(f)
			if header[0] == 0x1f && header[1] == 0x8b {
				reader, err = gzip.NewReader(f)
				if err != nil {
					return err
				}
			}

			skip, err := a.importCIF(reader)
			if err != nil {
				if skip {
					// Non fatal error so log it but don't kill the import
					log.Println(err)
				} else {
					return err
				}
			} else {
				fileCount++
			}
		}
	}

	// Do maintenance if in maintenance mode or we imported at least 1 CIF file
	if *(a.maintenance) || *(a.forceExpire) || fileCount > 0 {
		err := a.cleanup()
		if err != nil {
			return err
		}
	}

	if *(a.maintenance) || *(a.forceVacuum) || fileCount > 0 {
		err := a.vacuum()
		if err != nil {
			return err
		}

		err = a.cluster()
		if err != nil {
			return err
		}
	}

	if *(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum) {
		log.Println("Maintenance complete")
	} else {
		log.Println("Import complete")
	}

	return nil
}

func (c *CIFImporter) Update(f func(*sql.Tx) error) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
