package cif

func (c *CIF) parseTI( l string ) error {
  var t Tiploc = Tiploc{}
  i := 2
  i = parseStringTrim( l, i, 7, &t.Tiploc )
  i += 2
  i = parseInt( l, i, 6, &t.NLC )
  i = parseString( l, i, 1, &t.NLCCheck )
  i = parseStringTitle( l, i, 26, &t.Desc )
  i = parseInt( l, i, 5, &t.Stanox )
  i += 4
  i = parseStringTrim( l, i, 3, &t.CRS )
  i = parseStringTitle( l, i, 16, &t.NLCDesc )

  return c.put( c.tiploc, t.Tiploc, &t )
}
