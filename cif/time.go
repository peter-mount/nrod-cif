package cif

import (
  "fmt"
)

// Public Timetable time
// Note: 00:00 is not possible as in CIF that means no-time
type PublicTime struct {
  t int
}

func (t *PublicTime) String() string {
  if t.t <= 0 {
    return "     "
  }

  return fmt.Sprintf( "%02d:%02d", t.t/3600, (t.t/60)%60 )
}

func (t *PublicTime) Get() int {
  return t.t
}

func (t *PublicTime) Set( v int ) {
  t.t = v
}

func (t *PublicTime) IsSet() bool {
  return t.t<=0
}

// Working Timetable time
type WorkingTime struct {
  t int
}

func (t *WorkingTime) String() string {
  if t.t < 0 {
    return "        "
  }

  return fmt.Sprintf( "%02d:%02d:%02d", t.t/3600, (t.t/60)%60, t.t%60 )
}

func (t *WorkingTime) Get() int {
  return t.t
}

func (t *WorkingTime) Set( v int ) {
  t.t = v
}

func (t *WorkingTime) IsSet() bool {
  return t.t<0
}
