package cif

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
)

// TiplocHandler implements a net/http handler that implements a simple Rest service to retrieve Tiploc records.
// The handler must have {id} set in the path for this to work, where id would represent the Tiploc code.
//
// For example:
//
// router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /tiploc/MSTONEE to return JSON representing that station.
func (c *CIF) TiplocHandler( r *rest.Rest ) error {
  return c.db.View( func( tx *bolt.Tx ) error {
    tpl := r.Var( "id" )

    response := NewResponse()
    r.Value( response )

    if tiploc, exists := c.GetTiploc( tx, tpl ); exists {
      statistics.Incr( "tiploc.200" )
      r.Status( 200 )
      response.Status = 200
      response.AddTiploc( tiploc )
      tiploc.SetSelf( r )
      response.Self = tiploc.Self
    } else {
      statistics.Incr( "tiploc.404" )
      r.Status( 404 )
      response.Status = 404
      response.Message = tpl + " not found"
    }

    return nil
  } )
}
