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

func (t *PublicTime) String() string {
  if t.T <= 0 {
    return "     "
  }

  return fmt.Sprintf( "%02d:%02d", t.T/3600, (t.T/60)%60 )
}

func (t *PublicTime) Get() int {
  return t.T
}

func (t *PublicTime) Set( v int ) {
  t.T = v
}

func (t *PublicTime) IsSet() bool {
  return t.T<=0
}

// Working Timetable time
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

func (t *WorkingTime) String() string {
  if t.T < 0 {
    return "        "
  }

  return fmt.Sprintf( "%02d:%02d:%02d", t.T/3600, (t.T/60)%60, t.T%60 )
}

func (t *WorkingTime) Get() int {
  return t.T
}

func (t *WorkingTime) Set( v int ) {
  t.T = v
}

func (t *WorkingTime) IsSet() bool {
  return t.T<0
}
