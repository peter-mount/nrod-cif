package cifrest

import (
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/nrod-cif/cif"
  "fmt"
  "log"
  "strconv"
)

// TiplocHandler implements a net/http handler that implements a simple Rest service to retrieve Tiploc records.
// The handler must have {id} set in the path for this to work, where id would represent the Tiploc code.
//
// For example:
//
// router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /tiploc/MSTONEE to return JSON representing that station.
func (c *CIFRest) StanoxHandler( r *rest.Rest ) error {
  stanox, err := strconv.Atoi( r.Var( "id" ) )
  if err != nil {
    return err
  }

  tiplocs, err := c.cif.GetStanox( stanox )

  if err != nil {
    r.Status( 500 )
    log.Printf( "500: stanox %s = %s", stanox, err )
    return err
  }

  if tiplocs == nil || len( tiplocs ) == 0 {
    r.Status( 404 )
    return nil
  }

  resp := &cif.Response{
    Stanox: stanox,
    Self: r.Self( fmt.Sprintf( "/stanox/%d", stanox ) ),
  }
  resp.AddTiplocs( tiplocs )
  r.Status( 200 ).Value( resp )
  return nil
}
