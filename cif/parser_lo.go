package cif

import (
  "strings"
)

func (c *CIF) parseLO( l string ) error {
  s := c.curSchedule

  var loc *Location = &Location{}
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )
  loc.Tiploc = strings.Trim( loc.Location[0:8], " " )

  i = parseHHMMS( l, i, &loc.Wtd )
  i = parseHHMM( l, i, &loc.Ptd )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Line )

  i = parseStringTrim( l, i, 2, &loc.EngAllow )
  i = parseStringTrim( l, i, 2, &loc.PathAllow )

  i = parseActivity( l, i, &loc.Activity)

  i = parseStringTrim( l, i, 2, &loc.PerfAllow )

  s.appendLocation( loc )

  return nil
}
