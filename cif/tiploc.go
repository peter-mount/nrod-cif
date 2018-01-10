package cif

import (
  "encoding/json"
  "fmt"
  bolt "github.com/coreos/bbolt"
  "github.com/gorilla/mux"
  "github.com/peter-mount/golib/statistics"
  "log"
  "net/http"
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

func (c *CIF) parseTiplocInsert( l string ) error {
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

func (c *CIF) parseTiplocAmend( l string ) error {
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

func (c *CIF) parseTiplocDelete( l string ) error {
  var t Tiploc = Tiploc{}
  i := 2
  i = parseStringTrim( l, i, 7, &t.Tiploc )
  return c.tiploc.Delete( []byte( t.Tiploc ) )
}

func (c *CIF) GetTiploc( tx *bolt.Tx, t string ) ( *Tiploc, bool ) {
  var tiploc *Tiploc = &Tiploc{}

  if err := c.get( tx.Bucket( []byte("Tiploc") ), t, tiploc ); err != nil {
    return nil, false
  }

  return tiploc, true
}

func (c *CIF) TiplocHandler( w http.ResponseWriter, r *http.Request ) {
  var params = mux.Vars( r )

  tpl := params[ "id" ]

  tx, err := c.db.Begin(true)
  if err != nil {
    statistics.Incr( "tiploc.500" )
    log.Println( "Get Tiploc", tpl, err )
    w.WriteHeader( 500 )
    return
  }
  defer tx.Rollback()

  if tiploc, exists := c.GetTiploc( tx, tpl ); exists {

    if err := tx.Commit(); err != nil {
      statistics.Incr( "tiploc.500" )
      log.Println( "Get Tiploc", tpl, err )
      w.WriteHeader( 500 )
    } else {
      statistics.Incr( "tiploc.200" )
      w.WriteHeader( 200 )
      json.NewEncoder( w ).Encode( tiploc )
    }
  } else {
    statistics.Incr( "tiploc.404" )
    w.WriteHeader( 404 )
  }
}
