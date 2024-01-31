package cifimport

import (
	"log"
)

// cleanup removes expired entries from the database
func (c *CIFImporter) cleanup() error {
	log.Println("Removing historic associations")
	res, err := c.db.Exec("DELETE FROM timetable.assoc WHERE enddate < NOW()::DATE")
	if err != nil {
		return err
	}
	rc, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rc > 0 {
		log.Printf("Removed %d associations", rc)
	}

	log.Println("Removing historic schedules")
	res, err = c.db.Exec("DELETE FROM timetable.schedule WHERE enddate < NOW()::DATE")
	if err != nil {
		return err
	}
	rc, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if rc > 0 {
		log.Printf("Removed %d schedules", rc)
	}

	log.Println("Fixing tiploc crs codes")
	_, err = c.db.Exec("SELECT timetable.fixtiploccrs()")
	if err != nil {
		return err
	}

	return nil
}
