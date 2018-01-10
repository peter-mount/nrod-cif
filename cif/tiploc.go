package cif

import (
  "fmt"
)

type Tiploc struct {
  Tiploc    string
  NLC       int
  NLCCheck  string
  Desc      string
  Stanox    int
  CRS       string
  NLCDesc   string
}

func (t *Tiploc) String() string {
  return fmt.Sprintf(
    "Tiploc[%s, crs=%s, stanox=%05d, nlc=%d, desc=%s, nlcDesc=%s]",
    t.Tiploc,
    t.CRS,
    t.Stanox,
    t.NLC,
    t.Desc,
    t.NLCDesc )
}

func (c *CIF) parseTiplocInsert( l string ) {
  var t *Tiploc = &Tiploc{}
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

  c.tiploc[ t.Tiploc ] = t
}

func (c *CIF) parseTiplocAmend( l string ) {
  var t *Tiploc = &Tiploc{}
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
    c.tiploc[ t.Tiploc ] = t
  } else {
    delete( c.tiploc, t.Tiploc )
    c.tiploc[ newTiploc ] = t
  }

}

func (c *CIF) parseTiplocDelete( l string ) {
  var t *Tiploc = &Tiploc{}
  i := 2
  i = parseStringTrim( l, i, 7, &t.Tiploc )
  delete( c.tiploc, t.Tiploc )
}

func (c *CIF) GetTiploc( t string ) ( *Tiploc, bool ) {
  r, e := c.tiploc[ t ]
  return r, e
}

func (c *CIF) GetCRS( t string ) ( []*Tiploc, bool ) {
  r, e := c.crs[ t ]
  return r, e
}

func (c *CIF) GetStanox( s int ) ( []*Tiploc, bool ) {
  r, e := c.stanox[ s ]
  return r, e
}
