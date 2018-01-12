package cif

import (
  "encoding/json"
  bolt "github.com/coreos/bbolt"
  "github.com/gorilla/mux"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/statistics"
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
    codec.NewBinaryCodecFrom( v ).Read( tiploc )

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

    codec := codec.NewBinaryCodec()
    codec.WriteStringArray( ar )
    if codec.Error() != nil {
      return codec.Error()
    }

    if err := c.crs.Put( []byte( k ), codec.Bytes() ); err != nil {
      return err
    }
  }

  return nil
}

// GetCRS retrieves an array of Tiploc records for the CRS/3Alpha code of a station.
func (c *CIF) GetCRS( tx *bolt.Tx, crs string ) ( []*Tiploc, bool ) {

  b := tx.Bucket( []byte("Crs") ).Get( []byte( crs ) )
  if b == nil {
    return nil, false
  }

  var ar []string
  codec.NewBinaryCodecFrom( b ).ReadStringArray( &ar )

  if len( ar ) == 0 {
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

// CRSHandler implements a net/http handler that implements a simple Rest service to retrieve CRS/3Alpha records.
// The handler must have {id} set in the path for this to work, where id would represent the CRS code.
//
// For example:
//
// router.HandleFunc( "/crs/{id}", db.CRSHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /crs/MDE to return JSON representing that station.
func (c *CIF) CRSHandler( w http.ResponseWriter, r *http.Request ) {
  var params = mux.Vars( r )

  crs := params[ "id" ]

  if err := c.db.View( func( tx *bolt.Tx ) error {
    if ary, exists := c.GetCRS( tx, crs ); exists {
      statistics.Incr( "crs.200" )
      w.WriteHeader( 200 )
      json.NewEncoder( w ).Encode( ary )
    } else {
      statistics.Incr( "crs.404" )
      w.WriteHeader( 404 )
    }

    return nil
  }); err != nil {
    log.Println( "Get CRS", crs, err )
    statistics.Incr( "crs.500" )
    w.WriteHeader( 500 )
  }
}
