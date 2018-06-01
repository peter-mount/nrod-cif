package cif

import (
  "time"
)

type Association struct {
  MainUid         string    `json:"main"`
  AssocUid        string    `json:"assoc"`
  STPIndicator    string    `json:"stp"`
  StartDate       time.Time `json:"start"`
  EndDate         time.Time `json:"end"`
  AssocDays       string    `json:"dow"`
  Category        string    `json:"category"`
  DateInd         string    `json:"dateInd"`
  Tiploc          string    `json:"tpl"`
  BaseSuffix      string    `json:"baseSuffix,omitempty"`
  AssocSuffix     string    `json:"assocSuffix,omitempty"`
  AssocType       string    `json:"assocType"`
  // The CIF extract this entry is from
  DateOfExtract   time.Time `json:"date" xml:"date,attr"`
  // Self (generated on rest only)
  Self            string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}
