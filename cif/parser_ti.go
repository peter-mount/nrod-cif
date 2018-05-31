package cif

import (
  "log"
)

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

  return c.putTiploc( &t )
}

// Store/replace a tiploc only if the entry is newer than an existing one
func (c *CIF) putTiploc( t *Tiploc ) error {
  t.Update()

  // Link it to this CIF file
  t.DateOfExtract = c.importhd.DateOfExtract

  _, err := c.tx.Exec(
    "INSERT INTO timetable.tiploc " +
    "(tiploc, crs, stanox, name, nlc, nlccheck, nlcdesc, station, dateextract) " +
    "VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9)",
    t.Tiploc,
    t.CRS,
    t.Stanox,
    t.Desc,
    t.NLC,
    t.NLCCheck,
    t.NLCDesc,
    t.Station,
    t.DateOfExtract,
  )
  if err != nil {
    log.Printf( "Failed to insert tiploc %s", t.Tiploc )
    return err
  }

  return nil
}
