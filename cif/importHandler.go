package cif

import (
  "bufio"
  "encoding/json"
  "log"
  "net/http"
)

type SimpleResponse struct {
  Status  int
  Message string
}

func (c *CIF) ImportHandler( rw http.ResponseWriter, req *http.Request ) {
  log.Println( "CIF Import started" )
  scanner := bufio.NewScanner( req.Body )

  var result SimpleResponse = SimpleResponse{
    Status: 200,
    Message: "",
  }

  if err := c.parseFile( scanner ); err != nil {
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
