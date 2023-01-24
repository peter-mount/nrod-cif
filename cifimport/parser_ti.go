package cifimport

import (
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFImporter) parseTI(l string) error {
	var t cif.Tiploc = cif.Tiploc{}
	i := 2
	i = parseStringTrim(l, i, 7, &t.Tiploc)
	i += 2
	i = parseInt(l, i, 6, &t.NLC)
	i = parseString(l, i, 1, &t.NLCCheck)
	i = parseStringTitle(l, i, 26, &t.Name)
	i = parseInt(l, i, 5, &t.Stanox)
	i += 4
	i = parseStringTrim(l, i, 3, &t.CRS)
	i = parseStringTitle(l, i, 16, &t.NLCDesc)

	return c.putTiploc(&t)
}

// Store/replace a tiploc only if the entry is newer than an existing one
func (c *CIFImporter) putTiploc(t *cif.Tiploc) error {
	t.Update()

	// Link it to this CIF file
	t.DateOfExtract = c.importhd.DateOfExtract

	_, err := c.tx.Exec(
		"INSERT INTO timetable.tiploc "+
			"(tiploc, crs, stanox, name, nlc, nlccheck, nlcdesc, station, dateextract) "+
			"VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9) "+
			"ON CONFLICT ( id ) "+
			"DO UPDATE SET "+
			"crs = EXCLUDED.crs, "+
			"stanox = EXCLUDED.stanox, "+
			"name = EXCLUDED.name, "+
			"nlc = EXCLUDED.nlc, "+
			"nlccheck = EXCLUDED.nlccheck, "+
			"nlcdesc = EXCLUDED.nlcdesc, "+
			"station = EXCLUDED.station, "+
			"dateextract = EXCLUDED.dateextract ",
		t.Tiploc,
		t.CRS,
		t.Stanox,
		t.Name,
		t.NLC,
		t.NLCCheck,
		t.NLCDesc,
		t.Station,
		t.DateOfExtract,
	)
	if err != nil {
		log.Printf("Failed to insert tiploc %s", t.Tiploc)
		return err
	}

	return nil
}
