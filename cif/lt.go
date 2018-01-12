package cif

import (
  "strings"
)

func (c *CIF) parseLT( l string ) error {
  s := c.curSchedule

  var loc *Location = &Location{}
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )
  loc.Tiploc = strings.Trim( loc.Location[0:8], " " )

  i = parseHHMMS( l, i, &loc.Wta )

  i = parseHHMM( l, i, &loc.Pta )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Path )
  i = parseActivity( l, i, &loc.Activity )

  s.appendLocation( loc )

  return nil
}
