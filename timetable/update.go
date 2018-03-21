package timetable

import (
  bolt "github.com/coreos/bbolt"
  "cif"
  "log"
)

type index struct {
  cif *cif.CIF
  index map[string]map[int][]string
}

// Update updates the timetable
// c Cursor on the CIF database in view mode
func ( t * Timetable ) Update( c *cif.CIF, bucket *bolt.Bucket ) error {
  log.Println( "Updating timetable" )

  // Index of crs/hour/schedules
  index := &index{ cif: c, index: make( map[string]map[int][]string ) }

  if err := bucket.ForEach( index.process ); err != nil {
    log.Println( "Timetable update failed", err )
    return err
  }

  log.Println( "Timetable complete" )
  return nil
}

func (i *index) getSlot( c string, h int ) []string {
  m, exists := i.index[ c ]
  if !exists {
    m := make( map[int][]string )
    i.index[ c ] = m
  }

  return m[ h ]
}

func (i *index ) setSlot( c string, h int, s []string ) {
  i.index[ c ][ h ] = s
}
