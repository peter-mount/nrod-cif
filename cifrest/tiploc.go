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
func (c *CIFRest) TiplocHandler(r *rest.Rest) error {
  tpl := r.Var("id")

  tiploc, err := c.cif.GetTiploc(tpl)

  if tiploc == nil {
    r.Status(404)
    return nil
  }

  if err != nil {
    r.Status(500)
    log.Printf("500: tiploc %s = %s", tpl, err)
    return err
  }

  tiploc.Self = r.Self("/tiploc/" + tpl)
  r.Status(200).Value(tiploc)
  return nil
}
