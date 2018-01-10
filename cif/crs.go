package cif

import (
  "encoding/json"
  bolt "github.com/coreos/bbolt"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "sort"
)

func (c *CIF) cleanupCRS() error {
  log.Println( "Rebuilding CRS bucket" )

  // Clear the crs bucket
  if err := c.crs.ForEach( func( k, v []byte) error {
    return c.crs.Delete( k )
  }); err != nil {
    return err
  }

  // Refresh CRS map
  crs := make( map[string][]*Tiploc )

  if err := c.tiploc.ForEach( func( k, v []byte) error {
    var tiploc *Tiploc = &Tiploc{}
    if err := getInterface( v, tiploc ); err != nil {
      return err
    }

    if tiploc.CRS != "" {
      crs[ tiploc.CRS ] = append( crs[ tiploc.CRS ], tiploc )
    }

    return nil
  }); err != nil {
    return err
  }

  // Sort each crs slice by NLC, hopefully making the more accurate entry first
  // e.g. Look at VIC as an example
  for _, t := range crs {
    if len( t ) > 1 {
      sort.SliceStable( t, func( i, j int ) bool {
        return t[i].NLC < t[j].NLC
      })
    }
  }

  // Now persist
  for k, v := range crs {
    // Array of just Tiploc codes to save space
    var ar []string
    for _, t := range v {
      ar = append( ar, t.Tiploc )
    }

    c.put( c.crs, k, ar )
  }

  return nil
}

func (c *CIF) GetCRS( tx *bolt.Tx, crs string ) ( []*Tiploc, bool ) {

  var ar []string

  if err := c.get( tx.Bucket( []byte("Crs") ), crs, &ar ); err != nil {
    log.Println( err )
    return nil, false
  }

  var t []*Tiploc
  for _, k := range ar {
    if tiploc, exists := c.GetTiploc( tx, k ); exists {
      t = append( t, tiploc )
    }
  }

  return t, len( t ) > 0
}

func (c *CIF) CRSHandler( w http.ResponseWriter, r *http.Request ) {
  var params = mux.Vars( r )

  crs := params[ "id" ]

  tx, err := c.db.Begin(true)
  if err != nil {
    log.Println( "Get CRS", crs, err )
    w.WriteHeader( 500 )
    return
  }
  defer tx.Rollback()

  if ary, exists := c.GetCRS( tx, crs ); exists {

    if err := tx.Commit(); err != nil {
      log.Println( "Get CRS", crs, err )
      w.WriteHeader( 500 )
    } else {
      w.WriteHeader( 200 )
      json.NewEncoder( w ).Encode( ary )
    }
  } else {
    w.WriteHeader( 404 )
  }
}
