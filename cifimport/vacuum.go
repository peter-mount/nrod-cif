package cifimport

import (
	"context"
	"log"
)

// vacuum performs a vacuum full on the various tables
func (c *CIFImporter) vacuum(_ context.Context) error {
	log.Println("Compacting tiplocs")
	_, err := c.db.Exec("VACUUM FULL timetable.tiploc")
	if err != nil {
		return err
	}

	log.Println("Compacting associations")
	_, err = c.db.Exec("VACUUM FULL timetable.assoc")
	if err != nil {
		return err
	}

	log.Println("Compacting schedules")
	_, err = c.db.Exec("VACUUM FULL timetable.schedule")
	if err != nil {
		return err
	}
	_, err = c.db.Exec("VACUUM FULL timetable.schedule_json")
	if err != nil {
		return err
	}

	log.Println("Compacting station index")
	_, err = c.db.Exec("VACUUM FULL timetable.station")
	if err != nil {
		return err
	}

	return nil
}
