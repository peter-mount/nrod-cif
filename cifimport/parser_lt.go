package cifimport

import (
  "github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFImporter) parseLT( l string ) error {
  s := c.curSchedule

  var loc *cif.Location = &cif.Location{Type:"LT"}
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )

  i = parseHHMMS( l, i, &loc.Times.Wta )

  i = parseHHMM( l, i, &loc.Times.Pta )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Path )
  i = parseActivity( l, i, &loc.Activity )

  s.AppendLocation( loc )

  return nil
}
