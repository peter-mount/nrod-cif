// NR CIF file format

package cif

import (
  "fmt"
  "sort"
)

type CIF struct {
  // Copy of latest HD record
  Header    *HD
  // Map of Tiploc's
  tiploc    map[string]*Tiploc
  // Map of CRS codes to Tiplocs
  crs       map[string][]*Tiploc
  // Map of Stanox to Tiplocs
  stanox    map[int][]*Tiploc
  // Map of Schedules
  schedules map[string][]*Schedule
}

// Initialise a blank CIF
func (c *CIF ) Init() *CIF {
  c.Header = &HD{}
  c.tiploc = make( map[string]*Tiploc )
  c.crs = make( map[string][]*Tiploc )
  c.stanox = make( map[int][]*Tiploc )
  c.schedules = make( map[string][]*Schedule )
  return c
}

func (c *CIF) String() string {
  return fmt.Sprintf(
    "CIF %s Extracted %v Date Range %v - %v\ntiploc %d\ncrs %d\nschedules %d",
    c.Header.FileMainframeIdentity,
    c.Header.DateOfExtract.Format( HumanDateTime ),
    c.Header.UserStartDate.Format( HumanDate ),
    c.Header.UserEndDate.Format( HumanDate ),
    len( c.tiploc ),
    len( c.crs ),
    len( c.schedules ) )
}

func (c *CIF) cleanup() {
  c.cleanupStanox()
  c.cleanupCRS()
  c.cleanupSchedules()
}

func (c *CIF) cleanupStanox() {
  // Refresh stanox map
  if c.stanox == nil || len( c.stanox ) > 0 {
    c.stanox = make( map[int][]*Tiploc )
  }

  for _, t := range c.tiploc {
    if t.Stanox > 0 {
      c.stanox[ t.Stanox ] = append( c.stanox[ t.Stanox ], t )
    }
  }

  // Now for each stanox, if 1 entry has a crs then use that for all entries
  for _, s := range c.stanox {
    var crs string
    for _, t := range s {
      // Don't use X?? or Z?? CRS codes here
      if t.CRS != "" && !( t.CRS[0:1]=="X" || t.CRS[0:1]=="Z" ) {
        crs = t.CRS
      }
    }

    // Update to the new crs field
    if crs != "" {
      for _, t := range s {
        t.CRS = crs
      }
    }

    // Sort the slice by NLC, hopefully making the more accurate entry first
    if len( s ) > 1 {
      sort.SliceStable( s, func( i, j int ) bool {
        return s[i].NLC < s[j].NLC
      })
    }

  }

}

func (c *CIF) cleanupCRS() {
  // Refresh CRS map
  if c.crs == nil || len( c.crs ) > 0 {
    // Clear the CRS map
    c.crs = make( map[string][]*Tiploc )
  }

  for _, t := range c.tiploc {
    if t.CRS != "" {
      c.crs[ t.CRS ] = append( c.crs[ t.CRS ], t )
    }
  }

  // Sort each crs slice by NLC, hopefully making the more accurate entry first
  for _, t := range c.crs {
    if len( t ) > 1 {
      sort.SliceStable( t, func( i, j int ) bool {
        return t[i].NLC < t[j].NLC
      })
    }
  }

}


func (c *CIF) cleanupSchedules() {
  // Sort each schedule slice in start date & STP Indicator order, C, N, O & P
  for _, s := range c.schedules {
    if len( s ) > 1 {
      sort.SliceStable( s, func( i, j int ) bool {
        return s[i].RunsFrom.Before( s[j].RunsFrom ) && s[i].STPIndicator < s[i].STPIndicator
      })
    }
  }
}

// Returns all schedules for a train uid
func (c *CIF) GetSchedules( uid string ) []*Schedule {
  return c.schedules[ uid ]
}

func (c *CIF) addSchedule( s *Schedule ) {
  if ary, exists := c.schedules[ s.TrainUID ]; exists {
    // Check to see if we have a comparable entry. If so then replace it
    for i, e := range ary {
      if s.Equals( e ) {
        ary[ i ] = s
        return
      }
    }
  }

  c.schedules[ s.TrainUID ] = append( c.schedules[ s.TrainUID ], s )
}

func (c *CIF) deleteSchedule( s *Schedule ) {
  if ary, exists := c.schedules[ s.TrainUID ]; exists {
    var n []*Schedule
    for _, e := range ary {
      if !s.Equals( e ) {
        n = append( n, e )
      }
    }
    if len( n ) > 0 {
      c.schedules[ s.TrainUID ] = n
    } else {
      delete( c.schedules, s.TrainUID )
    }
  }
}
