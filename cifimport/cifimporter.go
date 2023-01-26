package cifimport

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/peter-mount/go-kernel/v2/db"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util/task"
	common "github.com/peter-mount/nrod-cif"
	"github.com/peter-mount/nrod-cif/cif"
	_ "gopkg.in/yaml.v2"
	"os"
)

type CIFImporter struct {
	dbService   *db.DBService `kernel:"inject"`
	worker      task.Queue    `kernel:"worker"`
	maintenance *bool         `kernel:"flag,m,Same as -expire -vacuum"`
	forceExpire *bool         `kernel:"flag,expire,Remove expired entries"`
	forceVacuum *bool         `kernel:"flag,vacuum,Vacuum and re-cluster database"`
	fileSource  *string       `kernel:"flag,files,files containing cif files to import"`
	db          *sql.DB
	header      *HD     // Last import HD record
	importhd    *HD     // Current import HD record
	tx          *sql.Tx // === Entries used during import only
	curSchedule *cif.Schedule
	update      bool
}

func (a *CIFImporter) PostInit() error {

	files, err := a.addCIFFilesForImport()
	if err != nil {
		return err
	}

	if *(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum) {
		if len(files) > 0 {
			return fmt.Errorf("CIF files not permitted in maintenance mode")
		}
	} else if len(files) == 0 {
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
			files = append(files, scanner.Text())
		}
		err = scanner.Err()
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func (a *CIFImporter) Start() error {
	log.Println(common.Version)

	files, err := a.addCIFFilesForImport()
	if err != nil {
		return err
	}

	a.db = a.dbService.GetDB()
	if a.db == nil {
		return fmt.Errorf("No database")
	}

	// Normal mode, cleanup & import CIF files
	if !(*(a.maintenance) || *(a.forceExpire) || *(a.forceVacuum)) {
		// Do a cleanup first as it will remove expired entries freeing up some space
		a.worker.AddPriorityTask(899, a.cleanup)

		for idx, file := range files {
			a.worker.AddPriorityTask(900+idx, task.Of(a.importCIFTask).WithValue("file", file))
		}
	}

	// Do maintenance if in maintenance mode or we imported at least 1 CIF file
	maintenance := len(files) > 0 || *(a.maintenance)
	if maintenance || *(a.forceExpire) {
		a.worker.AddPriorityTask(1001, a.cleanup)
	}

	if maintenance || *(a.forceVacuum) {
		a.worker.AddPriorityTask(1001, a.vacuum)
	}

	if maintenance || *(a.forceExpire) || *(a.forceVacuum) {
		a.worker.AddPriorityTask(1001, a.cluster)
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
