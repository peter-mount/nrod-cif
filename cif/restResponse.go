package cif

import (
  "encoding/xml"
  "github.com/peter-mount/golib/rest"
)

// Common struct used in forming all responses from rest endpoints.
// This makes the responses similar in nature and reduces the amount of
// redundant code
type Response struct {
  XMLName       xml.Name    `json:"-" xml:"response"`
  Message       string      `json:"message,omitempty" xml:"message,attr,omitempty"`
  Schedules  []*Schedule    `json:"schedules,omitempty" xml:"schedules>schedule,omitempty"`
  Tiploc       *TiplocMap   `json:"tiploc,omitempty" xml:"tiplocs>tiploc,omitempty"`
  Self          string      `json:"self,omitempty" xml:"self,attr,omitempty"`
}

func NewResponse() *Response {
  return &Response{}
}

func ( resp *Response ) SetSelf( r *rest.Rest, s string ) {
  resp.Self = r.Self( s )
  r.Status( 200 ).Value( resp )
}
