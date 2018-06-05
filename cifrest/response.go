package cifrest

import (
  "cif"
  "encoding/xml"
)

// Common struct used in forming all responses from rest endpoints.
// The only exception are those which return a single instance like Tiploc
// This makes the responses similar in nature and reduces the amount of
// redundant code
type Response struct {
  XMLName       xml.Name      `json:"-" xml:"response"`
  Crs           string        `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  Stanox        int           `json:"stanox,omitempty" xml:"stanox,attr,omitempty"`
  Message       string        `json:"message,omitempty" xml:"message,attr,omitempty"`
  Schedules  []*cif.Schedule  `json:"schedules,omitempty" xml:"schedules>schedule,omitempty"`
  Tiploc       *TiplocMap     `json:"tiploc,omitempty" xml:"tiplocs>tiploc,omitempty"`
  Self          string        `json:"self,omitempty" xml:"self,attr,omitempty"`
}
