package cifimport

import (
  "github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFImporter) parseLO( l string ) error {
  s := c.curSchedule

  var loc *cif.Location = &cif.Location{Type:"LO"}
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )

  i = parseHHMMS( l, i, &loc.Times.Wtd )
  i = parseHHMM( l, i, &loc.Times.Ptd )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Line )

  i = parseStringTrim( l, i, 2, &loc.EngAllow )
  i = parseStringTrim( l, i, 2, &loc.PathAllow )

  i = parseActivity( l, i, &loc.Activity)

  i = parseStringTrim( l, i, 2, &loc.PerfAllow )

  s.AppendLocation( loc )

  return nil
}
