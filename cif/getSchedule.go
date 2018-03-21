package cif

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/codec"
  "time"
)

// GetSchedule returns a specific Schedule's for a specific TrainUID, startDate and STPIndicator
// If no schedule exists for the required key then nil is returned
func (c *CIF) GetSchedule( tx *bolt.Tx, uid string, date time.Time, stp string ) *Schedule {
  key := []byte( uid + date.Format( Date ) + stp )
  s := &Schedule{}
  b := tx.Bucket( []byte( "Schedule" ) ).Get( key )
  dec := codec.NewBinaryCodecFrom( b )
  dec.Read( s )
  if uid == s.ID.TrainUID {
    return s
  }
  return nil
}
