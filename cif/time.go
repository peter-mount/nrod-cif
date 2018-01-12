package cif

import (
  "fmt"
  "github.com/peter-mount/golib/codec"
)

// Public Timetable time
// Note: 00:00 is not possible as in CIF that means no-time
type PublicTime struct {
  T int
}

func (t PublicTime) Write( c *codec.BinaryCodec ) {
  c.WriteInt32( int32( t.T ) )
}

func (t PublicTime) Read( c *codec.BinaryCodec ) {
  var i int32
  c.ReadInt32( &i )
  t.T = int(i)
}

// String returns a PublicTime in HH:MM format or 5 blank spaces if it's not set.
func (t *PublicTime) String() string {
  if t.T <= 0 {
    return "     "
  }

  return fmt.Sprintf( "%02d:%02d", t.T/3600, (t.T/60)%60 )
}

// Get returns the PublicTime in seconds of the day
func (t *PublicTime) Get() int {
  return t.T
}

// Set sets the PublicTime in seconds of the day
func (t *PublicTime) Set( v int ) {
  t.T = v
}

// IsSet returns true if the time is set
func (t *PublicTime) IsSet() bool {
  return t.T<=0
}

// Working Timetable time.
// WorkingTime is similar to PublciTime, except we can have seconds.
// In the Working Timetable, the seconds can be either 0 or 30.
type WorkingTime struct {
  T int
}

func (t WorkingTime) Write( c *codec.BinaryCodec ) {
  c.WriteInt32( int32( t.T ) )
}

func (t WorkingTime) Read( c *codec.BinaryCodec ) {
  var i int32
  c.ReadInt32( &i )
  t.T = int(i)
}

// String returns a PublicTime in HH:MM:SS format or 8 blank spaces if it's not set.
func (t *WorkingTime) String() string {
  if t.T < 0 {
    return "        "
  }

  return fmt.Sprintf( "%02d:%02d:%02d", t.T/3600, (t.T/60)%60, t.T%60 )
}

// Get returns the WorkingTime in seconds of the day
func (t *WorkingTime) Get() int {
  return t.T
}

// Set sets the WorkingTime in seconds of the day
func (t *WorkingTime) Set( v int ) {
  t.T = v
}

// IsSet returns true if the time is set
func (t *WorkingTime) IsSet() bool {
  return t.T<0
}
