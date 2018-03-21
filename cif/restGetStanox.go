package cif

import (
  bolt "github.com/coreos/bbolt"
  "fmt"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "strconv"
)

// StanoxHandler implements a net/http handler that implements a simple Rest service to retrieve stanox records.
// The handler must have {id} set in the path for this to work, where id would represent the CRS code.
//
// For example:
//
// router.HandleFunc( "/stanox/{id}", db.StanoxHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /stanox/89403 to return JSON representing that station.
func (c *CIF) StanoxHandler( r *rest.Rest ) error {
  stanox, err := strconv.Atoi( r.Var( "id" ) )
  if err != nil {
    return err
  }

  return c.db.View( func( tx *bolt.Tx ) error {
    response := NewResponse()
    r.Value( response )

    if ary, exists := c.GetStanox( tx, stanox ); exists {
      statistics.Incr( "stanox.200" )
      r.Status( 200 )
      response.Status = 200
      response.AddTiplocs( ary )
      response.TiplocsSetSelf( r )
      response.sortTiplocs()
      response.Self = r.Self( fmt.Sprintf( "/stanox/%s", stanox ) )
    } else {
      statistics.Incr( "stanox.404" )
      r.Status( 404 )
      response.Status = 404
      response.Message = fmt.Sprintf( "%s not found", stanox )
    }

    return nil
  })
}
