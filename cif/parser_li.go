package cif

func (c *CIF) parseLI( l string ) error {
  s := c.curSchedule

  var loc *Location = &Location{}
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )

  i = parseHHMMS( l, i, &loc.Times.Wta )
  i = parseHHMMS( l, i, &loc.Times.Wtd )
  i = parseHHMMS( l, i, &loc.Times.Wtp )

  i = parseHHMM( l, i, &loc.Times.Pta )
  i = parseHHMM( l, i, &loc.Times.Ptd )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Line )
  i = parseStringTrim( l, i, 3, &loc.Path )
  i = parseActivity( l, i, &loc.Activity)

  i = parseStringTrim( l, i, 2, &loc.EngAllow )
  i = parseStringTrim( l, i, 2, &loc.PathAllow )
  i = parseStringTrim( l, i, 2, &loc.PerfAllow )

  s.appendLocation( loc )

  return nil
}
