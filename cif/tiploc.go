package cif

import (
  "database/sql"
  "encoding/xml"
  "fmt"
  "time"
)

// Tiploc represents a location on the rail network.
// This can be either a station, a junction or a specific point along the line/
type Tiploc struct {
  XMLName         xml.Name  `json:"-" xml:"tiploc"`
  // Tiploc key for this location
  Tiploc          string    `json:"tiploc" xml:"tiploc,attr"`
  // Proper description for this location
  Desc            string    `json:"desc,omitempty" xml:"desc,attr,omitempty"`
  // CRS code, "" for none. Codes starting with X or Z are usually not stations.
  CRS             string    `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  // Stannox code, 0 means none
  Stanox          int       `json:"stanox,omitempty" xml:"stanox,attr,omitempty"`
  // NLC
  NLC             int       `json:"nlc" xml:"nlc,attr"`
  NLCCheck        string    `json:"nlcCheck" xml:"nlcCheck,attr"`
  // NLC description of the location
  NLCDesc         string    `json:"nlcDesc,omitempty" xml:"nlcDesc,attr,omitempty"`
  // True if this tiploc is a station
  Station         bool      `json:"station,omitempty" xml:"station,attr,omitempty"`
  // The unique ID of this tiploc
  ID              int64     `json:"id" xml:"id,attr"`
  // The CIF extract this entry is from
  DateOfExtract   time.Time `json:"date" xml:"date,attr"`
  // Self (generated on rest only)
  Self            string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

func (t *Tiploc) Update() {
  // Tiploc is a station IF it has a stanox, crs & crs not start with X or Z
  t.Station = t.Stanox > 0 &&t.CRS != "" && !(t.CRS[0] == 'X' || t.CRS[0] == 'Z')
}

func (t *Tiploc) Scan( row Scannable ) (bool,error) {
  var del bool

  err := row.Scan(
    &t.ID,
    &t.Tiploc,
    &t.CRS,
    &t.Stanox,
    &t.Desc,
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
    t.CRS=""
  }

  return false, err
}

// GetTiploc returns a single Tiploc or nil
func (c *CIF) GetTiploc( tpl string ) (*Tiploc, error) {
  row := c.db.QueryRow( "SELECT * FROM timetable.tiploc WHERE tiploc=$1", tpl )
  r := &Tiploc{}
  del, err := r.Scan( row )
  if del || err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }

  return r, nil
}

// GetTiplocQuery returns 0 or more tiplocs based on a query
// join is an optional join clause. Alias t is that of the tiploc table.
// where is the where clause
// args are the arguments needed by the where clause
// returns a slice or nil if no results matched
func (c *CIF) GetTiplocQuery( join string, where string, args ...interface{} ) ([]*Tiploc, error) {
  rows, err := c.db.Query( "SELECT t.* FROM timetable.tiploc t " + join + " WHERE " + where, args... )
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var tiploc *Tiploc
  var res []*Tiploc
  for rows.Next() {
    if tiploc == nil {
      tiploc = &Tiploc{}
    }

    del, err := tiploc.Scan( rows )
    if err != nil {
      return nil, err
    }

    if !del {
      res = append( res, tiploc )
      tiploc = nil
    }
  }

  return res, nil
}

// GetStanox returns 0 or more tiplocs with the same stanox
func (c *CIF) GetStanox( stanox int ) ([]*Tiploc, error) {
  rows, err := c.GetTiplocQuery( "", "stanox=$1", stanox )
  return rows, err
}

// GetCRS returns all tiplocs who's stanox matches the tiploc with a CRS value.
//
// Note this performs a join to an entry with a crs value because in CIF a
// single tiploc will have a crs however for some schedules the schedules are
// associated with another without a crs.
//
// For example, CRS VIC for London Victoria matches tiploc VICTRIA. However
// that tiploc has no schedules. They are associated with VICTRIC and VICTRIE
// so the join will allow us to get those schedules.
//
// The downside is for some stations like Victoria, theres other tiplocs as well
// some representing platforms.
func (c *CIF) GetCRS( crs string ) ([]*Tiploc, error) {
  rows, err := c.GetTiplocQuery(
    "INNER JOIN timetable.tiploc c ON t.stanox=c.stanox",
    "c.crs=$1 AND c.stanox IS NOT NULL AND c.stanox > 0",
    crs )
  return rows, err
}

// GetTiplocs returns a slice based on the supplied slice of tiploc names
func (c *CIF) GetTiplocs( tiplocs []string ) ([]*Tiploc, error) {
  if len( tiplocs ) == 0 {
    return nil, nil
  }

  where := "tiploc IN ("
  args := []interface{}{}
  for i, arg := range tiplocs {
    if i>0 {
      where = where + ","
    }
    where = fmt.Sprintf( "%s$%d", where, i+1 )
    args = append( args, arg )
  }
  where = where + ")"

  rows, err := c.GetTiplocQuery( "", where, args... )
  return rows, err
}
