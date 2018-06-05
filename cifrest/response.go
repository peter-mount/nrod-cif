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
  // Schedule primary key
  TrainUID      string        `json:"uid,omitempty" xml:"uid,attr,omitempty"`
  Date          string        `json:"date,omitempty" xml:"date,attr,omitempty"`
  STPIndicator  string        `json:"stp,omitempty" xml:"stp,attr,omitempty"`
  //Message       string        `json:"message,omitempty" xml:"message,attr,omitempty"`
  // Schedules if any
  Schedules  []*cif.Schedule  `json:"schedules,omitempty" xml:"schedules>schedule,omitempty"`
  // Map of tiplocs in result or in schedules
  Tiploc       *TiplocMap     `json:"tiploc,omitempty" xml:"tiplocs>tiploc,omitempty"`
  // uri that could replicate this result (updates dependent)
  Self          string        `json:"self,omitempty" xml:"self,attr,omitempty"`
}
