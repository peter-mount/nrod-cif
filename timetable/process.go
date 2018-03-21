package timetable

import (
  "cif"
  "github.com/peter-mount/golib/codec"
)

func (i *index) process( k, b []byte ) error {
  s := &cif.Schedule{}
  dec := codec.NewBinaryCodecFrom( b )
  dec.Read( s )

  // Don't use Cancellations but TODO check we handle these in final results
  if s.ID.STPIndicator != "C" {

    key := s.Key()

    // Add locations that are stations to the map
    for _, l := range s.Locations {
      tpl := i.getTiploc( l.Tiploc )
      if tpl != nil && tpl.Station {
        hr := l.Times.Time.Get() / 3600
        i.getSlot( tpl.CRS, hr )[ key ] = nil
      }
    }

  }
  return nil
}
