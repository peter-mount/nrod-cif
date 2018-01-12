package cif

import (
  "encoding/json"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "strconv"
)

// Public Timetable time
// Note: 00:00 is not possible as in CIF that means no-time
type PublicTime struct {
  t int
}

func (t *PublicTime) Write( c *codec.BinaryCodec ) {
  c.WriteInt32( int32( t.t ) )
}

func (t *PublicTime) Read( c *codec.BinaryCodec ) {
  var i int32
  c.ReadInt32( &i )
  t.t = int(i)
}

func (t *PublicTime) MarshalJSON() ( []byte, error ) {
  if t.t <= 0 {
    return json.Marshal( nil )
  }
  return json.Marshal( t.String() )
}

func (t *PublicTime) UnmarshalJSON( data []byte ) error {
  var aux *string
  if err := json.Unmarshal( data, &aux ); err != nil {
    return err
  }
  if aux == nil {
    t.t = -1
  } else {
      a, _ := strconv.Atoi( (*aux)[0:2] )
      b, _ := strconv.Atoi( (*aux)[3:5] )
      t.t = (a *3600) + (b * 60)
  }
  return nil
}

// String returns a PublicTime in HH:MM format or 5 blank spaces if it's not set.
func (t *PublicTime) String() string {
  if t.t <= 0 {
    return "     "
  }

  return fmt.Sprintf( "%02d:%02d", t.t/3600, (t.t/60)%60 )
}

// Get returns the PublicTime in seconds of the day
func (t *PublicTime) Get() int {
  return t.t
}

// Set sets the PublicTime in seconds of the day
func (t *PublicTime) Set( v int ) {
  t.t = v
}

// IsSet returns true if the time is set
func (t *PublicTime) IsSet() bool {
  return t.t<=0
}

// Working Timetable time.
// WorkingTime is similar to PublciTime, except we can have seconds.
// In the Working Timetable, the seconds can be either 0 or 30.
type WorkingTime struct {
  t int
}

func (t *WorkingTime) Write( c *codec.BinaryCodec ) {
  c.WriteInt32( int32( t.t ) )
}

func (t *WorkingTime) Read( c *codec.BinaryCodec ) {
  var i int32
  c.ReadInt32( &i )
  t.t = int(i)
}

func (t *WorkingTime) MarshalJSON() ( []byte, error ) {
  if t.t < 0 {
    return json.Marshal( nil )
  }
  return json.Marshal( t.String() )
}

func (t *WorkingTime) UnmarshalJSON( data []byte ) error {
  var aux *string
  if err := json.Unmarshal( data, &aux ); err != nil {
    return err
  }
  if aux == nil {
    t.t = -1
  } else {
    a, _ := strconv.Atoi( (*aux)[0:2] )
    b, _ := strconv.Atoi( (*aux)[3:5] )
    c, _ := strconv.Atoi( (*aux)[6:8] )
    t.t = (a *3600) + (b * 60) + c
  }
  return nil
}

// String returns a PublicTime in HH:MM:SS format or 8 blank spaces if it's not set.
func (t *WorkingTime) String() string {
  if t.t < 0 {
    return "        "
  }

  return fmt.Sprintf( "%02d:%02d:%02d", t.t/3600, (t.t/60)%60, t.t%60 )
}

// Get returns the WorkingTime in seconds of the day
func (t *WorkingTime) Get() int {
  return t.t
}

// Set sets the WorkingTime in seconds of the day
func (t *WorkingTime) Set( v int ) {
  t.t = v
}

// IsSet returns true if the time is set
func (t *WorkingTime) IsSet() bool {
  return t.t<0
}
