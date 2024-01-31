package cifimport

import (
	"github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFImporter) parseTA(l string) error {
	var t cif.Tiploc = cif.Tiploc{}
	i := 2
	i = parseStringTrim(l, i, 7, &t.Tiploc)
	i += 2
	i = parseInt(l, i, 6, &t.NLC)
	i = parseStringTrim(l, i, 1, &t.NLCCheck)
	i = parseStringTrim(l, i, 26, &t.Name)
	i = parseInt(l, i, 5, &t.Stanox)
	i += 4
	i = parseStringTrim(l, i, 3, &t.CRS)
	i = parseStringTrim(l, i, 16, &t.NLCDesc)

	var newTiploc string
	i = parseStringTrim(l, i, 7, &newTiploc)

	if newTiploc != "" {

		_, err := c.tx.Exec("DELETE FROM timetable.tiploc WHERE tiploc = $1", t.Tiploc)
		if err != nil {
			return err
		}

		// Update the "new" entry to the new name
		t.Tiploc = newTiploc
	}

	// persist
	return c.putTiploc(&t)
}
