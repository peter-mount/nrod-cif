package cifimport

import (
	"context"
	"github.com/peter-mount/go-kernel/v2/log"
)

func (c *CIFImporter) cluster(_ context.Context) error {
	// These must be done outside of a transaction!

	// tiplocs clustered by stanox as we use that most often when searching
	// against crs & do the search on all tiplocs with the same stanox
	log.Println("Clustering tiplocs")
	_, err := c.db.Exec("CLUSTER timetable.tiploc USING tiploc_cluster")
	if err != nil {
		return err
	}

	// Cluster associations
	log.Println("Clustering associations")
	_, err = c.db.Exec("CLUSTER timetable.assoc USING assoc_cluster")
	if err != nil {
		return err
	}

	// Cluster schedules on their uid so we have all related schedules together.
	log.Println("Clustering schedules")
	_, err = c.db.Exec("CLUSTER timetable.schedule USING schedule_uid")
	if err != nil {
		return err
	}

	log.Println("Clustering station index")
	_, err = c.db.Exec("CLUSTER timetable.station USING station_tdt")
	return err
}
