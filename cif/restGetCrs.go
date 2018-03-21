package cif

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
)

// CRSHandler implements a net/http handler that implements a simple Rest service to retrieve CRS/3Alpha records.
// The handler must have {id} set in the path for this to work, where id would represent the CRS code.
//
// For example:
//
// router.HandleFunc( "/crs/{id}", db.CRSHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /crs/MDE to return JSON representing that station.
func (c *CIF) CRSHandler( r *rest.Rest ) error {
  return c.db.View( func( tx *bolt.Tx ) error {
    crs := r.Var( "id" )

    response := NewResponse()
    r.Value( response )

    if ary, exists := c.GetCRS( tx, crs ); exists {
      statistics.Incr( "crs.200" )
      r.Status( 200 )
      response.Status = 200
      response.AddTiplocs( ary )
      response.TiplocsSetSelf( r )
      response.sortTiplocs()
      response.Self = r.Self( "/crs/" + crs )
      // Set tiploc selfs
      for _, t := range ary {
        t.SetSelf( r )
      }
    } else {
      statistics.Incr( "crs.404" )
      r.Status( 404 )
      response.Status = 404
      response.Message = crs + " not found"
    }

    return nil
  })
}
