package cifimport

import (
	"github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFImporter) parseTD(l string) error {
	var t cif.Tiploc = cif.Tiploc{}
	i := 2
	i = parseStringTrim(l, i, 7, &t.Tiploc)
	return nil //c.tiploc.Delete( []byte( t.Tiploc ) )
}
