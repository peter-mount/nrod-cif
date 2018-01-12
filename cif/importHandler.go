package cif

import (
  "encoding/json"
  "log"
  "net/http"
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
func (c *CIF) ImportHandler( rw http.ResponseWriter, req *http.Request ) {
  log.Println( "CIF Import started" )

  var result SimpleResponse = SimpleResponse{
    Status: 200,
    Message: "",
  }

  if err := c.ImportCIF( req.Body ); err != nil {
    log.Printf( "CIF Import: %+v", err )
    result.Status = 500
    result.Message = err.Error()
  } else {
    log.Println( "CIF Import: completed" )
    result.Status = 200
    result.Message = "success"
  }

  rw.WriteHeader( result.Status )
  json.NewEncoder( rw ).Encode( result )
}
