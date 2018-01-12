package cif

import (
  "github.com/peter-mount/golib/codec"
)

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

  if newTiploc != "" {
    // Remove the old entry only if it's older than the current CIF file
    b := c.tiploc.Get( []byte( t.Tiploc ) )

    var ot Tiploc
    if( b != nil ) {
      codec.NewBinaryCodecFrom( b ).Read( &ot )
    }

    if t.Tiploc == ot.Tiploc && c.importhd.DateOfExtract.After( ot.DateOfExtract) {
      if err := c.tiploc.Delete( []byte( t.Tiploc ) ); err != nil {
        return err
      }
    }

    // Update the "new" entry to the new name
    t.Tiploc = newTiploc
  }

  // persist
  return c.putTiploc( &t )
}
