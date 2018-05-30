package cif

import (
  "github.com/peter-mount/golib/rest"
  "log"
)

type SimpleResponse struct {
  Status  int
  Message string
}

// ImportHandler implements a net/http handler that implements a Rest service to
// import an uncompressed CIF file from NetworkRail and import it into the cif
// database.
//
// For example:
//
// router.HandleFunc( "/importCIF", db.ImportHandler ).Methods( "POST" )
//
// Will define the path /importCIF to accept HTTP POST requests. You can then
// submit a cif file to this endpoint to import a CIF file.
//
// Example: To perform a full import, replacing all cif data in the database
//
// curl -X POST --data-binary @toc-full.CIF http://localhost:8081/importCIF
//
// To perform an update then simply submit an update cif file:
//
// curl -X POST --data-binary @toc-update-sun.CIF http://localhost:8081/importCIF
//
// BUG(peter-mount): The Rest service provided by ImportHandler is currently
// unprotected so anyone can perform an import.
// We need to provide some means of simple authentication to this handler.
func (c *CIF) ImportHandler( r *rest.Rest ) error {
  log.Println( "CIF Import started" )

  if err := c.ImportCIF( r.Request().Body ); err != nil {
    log.Printf( "CIF Import: %+v", err )
    return err
  }

  log.Println( "CIF Import: completed" )
  r.Status( 200 ).
    Value( &SimpleResponse{ Status: 200, Message: "success"} )

  return nil
}
