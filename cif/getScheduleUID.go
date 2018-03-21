package cif

import (
  bolt "github.com/coreos/bbolt"
  "bytes"
  "github.com/peter-mount/golib/codec"
)

// GetSchedulesByUID returns all Schedule's for a specific TrainUID.
// If no schedules exist for the required TrainUID then the returned slice is empty.
func (c *CIF) GetSchedulesByUID( tx *bolt.Tx, uid string ) []*Schedule {
  var ar []*Schedule

  b := tx.Bucket( []byte( "Schedule" ) ).Cursor()
  prefix := []byte( uid )

  for k, v := b.Seek( prefix ); k != nil && bytes.Compare( k[:len(prefix)], prefix ) == 0; k, v = b.Next() {
    s := &Schedule{}
    dec := codec.NewBinaryCodecFrom( v )
    dec.Read( s )
    if uid == s.TrainUID {
      ar = append( ar, s )
    }
  }

  return ar
}
