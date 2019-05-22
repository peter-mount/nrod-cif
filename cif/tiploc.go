package cif

import (
  "time"
)

// Tiploc represents a location on the rail network.
// This can be either a station, a junction or a specific point along the line/
type Tiploc struct {
  ID            int64     `json:"id" xml:"id,attr"`
  Tiploc        string    `json:"tiploc" xml:"tiploc,attr"`
  CRS           string    `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  Stanox        int       `json:"stanox,omitempty" xml:"stanox,attr,omitempty"`
  Name          string    `json:"name,omitempty" xml:"name,attr,omitempty"`
  NLC           int       `json:"nlc" xml:"nlc,attr"`
  NLCCheck      string    `json:"nlcCheck" xml:"nlcCheck,attr"`
  NLCDesc       string    `json:"nlcDesc,omitempty" xml:"nlcDesc,attr,omitempty"`
  Station       bool      `json:"station,omitempty" xml:"station,attr,omitempty"`
  Deleted       bool      `json:"deleted" xml:"deleted,attr,omitempty"`
  DateOfExtract time.Time `json:"date" xml:"date,attr"`
  Self          string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

func (t *Tiploc) Update() {
  // Tiploc is a station IF it has a stanox, crs & crs not start with X or Z
  t.Station = t.Stanox > 0 && t.CRS != "" && !(t.CRS[0] == 'X' || t.CRS[0] == 'Z')
}

func (t *Tiploc) Scan(row Scannable) (bool, error) {
  var del bool

  err := row.Scan(
    &t.ID,
    &t.Tiploc,
    &t.CRS,
    &t.Stanox,
    &t.Name,
    &t.NLC,
    &t.NLCCheck,
    &t.NLCDesc,
    &t.Station,
    &del,
    &t.DateOfExtract,
  )

  if err != nil {
    return false, err
  }

  if del {
    return true, nil
  }

  // Temporary fix until db uses nulls - bug in parser
  if t.CRS == "   " {
    t.CRS = ""
  }

  return false, err
}
