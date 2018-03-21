package timetable

import (
  bolt "github.com/coreos/bbolt"
  "cif"
  "log"
)

type index struct {
  cif *cif.CIF
  bucket *bolt.Bucket
  // Index of schedules keyed by crs/hour
  index map[string]map[int]map[string]interface{}
  // index of tiplocs, used to determine if we want to know about an entry
  tpl map[string]*cif.Tiploc
}

// Update updates the timetable
// c Cursor on the CIF database in view mode
func ( t * Timetable ) Update( c *cif.CIF, bucket *bolt.Bucket ) error {
  log.Println( "Updating timetable" )

  // Index of crs/hour/schedules
  index := &index{
    cif: c,
    bucket: bucket,
    index: make( map[string]map[int]map[string]interface{} ),
    tpl: make( map[string]*cif.Tiploc ),
  }

  if err := bucket.ForEach( index.process ); err != nil {
    log.Println( "Timetable update failed", err )
    return err
  }

  log.Println( "Timetable complete" )
  return nil
}

func (i *index) getSlot( c string, h int ) map[string]interface{} {
  m, exists := i.index[ c ]
  if !exists {
    m = make( map[int]map[string]interface{} )
    i.index[ c ] = m
  }

  s, exists := m[ h ]
  if !exists {
    s = make( map[string]interface{} )
    m[ h ] = s
  }

  return s
}

func (i *index) getTiploc( tpl string ) *cif.Tiploc {
  if t, exists := i.tpl[ tpl ]; exists {
    return t
  }

  if t, exists := i.cif.GetTiploc( i.bucket.Tx(), tpl ); exists {
    i.tpl[ tpl ] = t
    return t
  }

  return nil
}
