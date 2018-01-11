package cif

func (c *CIF) parseTA( l string ) error {
  var t Tiploc = Tiploc{}
  i := 2
  i = parseStringTrim( l, i, 7, &t.Tiploc )
  i += 2
  i = parseInt( l, i, 6, &t.NLC )
  i = parseStringTrim( l, i, 1, &t.NLCCheck )
  i = parseStringTrim( l, i, 26, &t.Desc )
  i = parseInt( l, i, 5, &t.Stanox )
  i += 4
  i = parseStringTrim( l, i, 3, &t.CRS )
  i = parseStringTrim( l, i, 16, &t.NLCDesc )

  var newTiploc string
  i = parseStringTrim( l, i, 7, &newTiploc )

  if newTiploc == "" {
    return c.put( c.tiploc, t.Tiploc, &t )
  } else {
    // Remove the old entry
    if err := c.tiploc.Delete( []byte( t.Tiploc ) ); err != nil {
      return err
    }

    // Update and store as the new entry
    t.Tiploc = newTiploc
    return c.put( c.tiploc, newTiploc, &t )
  }

}
