package cifrest

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

// TiplocHandler implements a net/http handler that implements a simple Rest service to retrieve Tiploc records.
// The handler must have {id} set in the path for this to work, where id would represent the Tiploc code.
//
// For example:
//
// router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /tiploc/MSTONEE to return JSON representing that station.
func (c *CIFRest) CRSHandler( r *rest.Rest ) error {
  crs := r.Var( "id" )

  tiplocs, err := c.cif.GetCRS( crs )

  if err != nil {
    r.Status( 500 )
    log.Printf( "500: crs %s = %s", crs, err )
    return err
  }

  if tiplocs == nil || len( tiplocs ) == 0 {
    r.Status( 404 )
    return nil
  }

  resp := &Response{
    Crs: crs,
    Self: r.Self( "/crs/" + crs ),
  }
  resp.AddTiplocs( tiplocs )
  r.Status( 200 ).Value( resp )
  return nil
}
