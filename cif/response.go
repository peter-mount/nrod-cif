package cif

import (
  "encoding/xml"
)

// Common struct used in forming all responses from rest endpoints.
// This makes the responses similar in nature and reduces the amount of
// redundant code
type Response struct {
  XMLName       xml.Name    `json:"-" xml:"response"`
  Status        int         `json:"status,omitempty" xml:"status,attr,omitempty"`
  Message       string      `json:"message,omitempty" xml:"message,attr,omitempty"`
  Schedules  []*Schedule    `json:"schedules,omitempty" xml:"schedules>schedule,omitempty"`
  //Tiploc     []*Tiploc              `json:"tiploc,omitempty" xml:"tiplocs>tiploc,omitempty"`
  Tiploc       *TiplocMap   `json:"tiploc,omitempty" xml:"tiplocs>tiploc,omitempty"`
  Self          string      `json:"self" xml:"self,attr,omitempty"`
}

func NewResponse() *Response {
  return &Response{}
}
