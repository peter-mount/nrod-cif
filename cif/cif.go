// NR CIF file format

package cif

import (
  bolt "github.com/coreos/bbolt"
  "errors"
  "fmt"
  "sort"
)

type CIF struct {
  db         *bolt.DB
  tx         *bolt.Tx
  // Copy of latest HD record
  Header     *HD
  // Map of Tiploc's
  //tiploc      map[string]*Tiploc
  tiploc     *bolt.Bucket
  // Map of CRS codes to Tiplocs
  crs        *bolt.Bucket
  // Map of Stanox to Tiplocs
  stanox     *bolt.Bucket
  // Map of Schedules
  schedules   map[string][]*Schedule
}

// Initialise a blank CIF
func (c *CIF ) Init( db *bolt.DB ) *CIF {
  c.db = db
  c.Header = &HD{}
  c.schedules = make( map[string][]*Schedule )
  return c
}

func (c *CIF) get( b *bolt.Bucket, k string, i interface{} ) error {
  bar := b.Get( []byte(k) )
  if bar != nil {
    return getInterface( bar, i )
  }
  return errors.New( k + " Not found")
}

func (c *CIF) put( b *bolt.Bucket, k string, i interface{} ) error {
  if bar, err := getBytes( i ); err != nil {
    return err
  } else {
    return b.Put( []byte(k), bar )
  }
}

func (c *CIF) resetDB() error {
  return c.tiploc.ForEach( func( k, v []byte) error {
    return c.tiploc.Delete( k )
  })
}

func (c *CIF) String() string {
  return fmt.Sprintf(
    "CIF %s Extracted %v Date Range %v - %v Update %s",
    c.Header.FileMainframeIdentity,
    c.Header.DateOfExtract.Format( HumanDateTime ),
    c.Header.UserStartDate.Format( HumanDate ),
    c.Header.UserEndDate.Format( HumanDate ),
    c.Header.Update )
}

func (c *CIF) Rebuild( tx *bolt.Tx ) error {

  c.tiploc = tx.Bucket( []byte("Tiploc") )
  c.crs = tx.Bucket( []byte("Crs") )
  c.stanox = tx.Bucket( []byte("Stanox") )

  if err := c.cleanupStanox(); err != nil {
    return err
  }

  if err := c.cleanupCRS(); err != nil {
    return err
  }

  //c.cleanupSchedules()

  return nil
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
