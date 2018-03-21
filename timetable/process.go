package timetable

import (
  "cif"
  "github.com/peter-mount/golib/codec"
)

func (i *index) process( k, b []byte ) error {
  s := &cif.Schedule{}
  dec := codec.NewBinaryCodecFrom( b )
  dec.Read( s )

//  for _, l := range s.Locations {
//
//  }

  return nil
}
