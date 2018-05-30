package cif

func (c *CIF) parseTD( l string ) error {
  var t Tiploc = Tiploc{}
  i := 2
  i = parseStringTrim( l, i, 7, &t.Tiploc )
  return nil //c.tiploc.Delete( []byte( t.Tiploc ) )
}
