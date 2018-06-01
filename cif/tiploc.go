package cif

import (
  "time"
)

// Tiploc represents a location on the rail network.
// This can be either a station, a junction or a specific point along the line/
type Tiploc struct {
  //XMLName         xml.Name  `xml:"tiploc"`
  // Tiploc key for this location
  Tiploc          string    `json:"tiploc" xml:"tiploc,attr"`
  NLC             int       `json:"nlc" xml:"nlc,attr"`
  NLCCheck        string    `json:"nlcCheck" xml:"nlcCheck,attr"`
  // Proper description for this location
  Desc            string    `json:"desc,omitempty" xml:"desc,attr,omitempty"`
  // Stannox code, 0 means none
  Stanox          int       `json:"stanox,omitempty" xml:"stanox,attr,omitempty"`
  // CRS code, "" for none. Codes starting with X or Z are usually not stations.
  CRS             string    `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  // NLC description of the location
  NLCDesc         string    `json:"nlcDesc,omitempty" xml:"nlcDesc,attr,omitempty"`
  // True if this tiploc is a station
  Station         bool      `json:"station"`
  // The CIF extract this entry is from
  DateOfExtract   time.Time `json:"date" xml:"date,attr"`
  // Self (generated on rest only)
  Self            string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

func (t *Tiploc) Update() {
  // Tiploc is a station IF it has a stanox, crs & crs not start with X or Z
  t.Station = t.Stanox > 0 &&t.CRS != "" && !(t.CRS[0] == 'X' || t.CRS[0] == 'Z')
}
