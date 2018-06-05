package db

import (
  "database/sql"
  "github.com/peter-mount/nrod-cif/cif"
  "fmt"
)

// GetTiploc returns a single Tiploc or nil
func (c *CIF) GetTiploc( tpl string ) (*cif.Tiploc, error) {
  row := c.db.QueryRow( "SELECT * FROM timetable.tiploc WHERE tiploc=$1", tpl )
  r := &cif.Tiploc{}
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
func (c *CIF) GetTiplocQuery( join string, where string, args ...interface{} ) ([]*cif.Tiploc, error) {
  rows, err := c.db.Query( "SELECT t.* FROM timetable.tiploc t " + join + " WHERE " + where, args... )
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var tiploc *cif.Tiploc
  var res []*cif.Tiploc
  for rows.Next() {
    if tiploc == nil {
      tiploc = &cif.Tiploc{}
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
func (c *CIF) GetStanox( stanox int ) ([]*cif.Tiploc, error) {
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
func (c *CIF) GetCRS( crs string ) ([]*cif.Tiploc, error) {
  rows, err := c.GetTiplocQuery(
    "INNER JOIN timetable.tiploc c ON t.stanox=c.stanox",
    "c.crs=$1 AND c.stanox IS NOT NULL AND c.stanox > 0",
    crs )
  return rows, err
}

// GetTiplocs returns a slice based on the supplied slice of tiploc names
func (c *CIF) GetTiplocs( tiplocs []string ) ([]*cif.Tiploc, error) {
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
