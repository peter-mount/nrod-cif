package cif

import (
  "github.com/peter-mount/golib/codec"
)

// A representation of a location within a schedule.
// There are three types of location, defined by the Id field:
//
// "LO" Origin, always the first location in a schedule
//
// "LI" Intermediate: A stop or pass along the route
//
// "LT" Destination: always the last lcoation in a schedule
//
// For most purposes you would be interested in the Tiploc, Pta, Ptd and Platform
// fields.
//
// Tiploc is the name of this location.
//
// Pta & Ptd are the public timetable times, i.e. what is published to the general public.
//
// Pta is the arrival time and is valid for LI & LT entries only.
//
// Ptd is the departue time and is valid for LO & LI entries only.
//
// If either are not set then the train is not scheduled to stop at this location.
//
// Wta, Wtd & Wtp are the working timetable, i.e. the actual timetable the
// service runs to. Wta & Wtd are like Pta & Ptd but Wtp means the time the train
// is scheduled to pass a location. If Wtp is set then Pta, Ptd, Wta & Wtp will
// not be set.
type Location struct {
  // Type of location:
  Id          string        `json:"-"`
  // Location including Suffix (for circular routes)
  // This is guaranteed to be unique per schedule, although for most purposes
  // like display you would use Tiploc
  Location    string        `json:"-"`
  // Tiploc of this location. For some schedules like circular routes this can
  // appear more than once in a schedule.
  Tiploc      string
  // Public Timetable
  Pta        *PublicTime    `json:",omitempty"`
  Ptd        *PublicTime    `json:",omitempty"`
  // Working Timetable
  Wta        *WorkingTime   `json:",omitempty"`
  Wtd        *WorkingTime   `json:",omitempty"`
  Wtp        *WorkingTime   `json:",omitempty"`
  // Platform
  Platform    string        `json:",omitempty"`
  // Activity up to 6 codes
  Activity  []string        `json:",omitempty"`
  // The Line the train will take
  Line        string        `json:",omitempty"`
  // The Path the train will take
  Path        string        `json:",omitempty"`
  // Allowances at this location
  EngAllow    string        `json:",omitempty"`
  PathAllow   string        `json:",omitempty"`
  PerfAllow   string        `json:",omitempty"`
}

// BinaryCodec writer
func (l *Location) Write( c *codec.BinaryCodec ) {
  c.WriteString( l.Id ).
    WriteString( l.Location ).
    WriteString( l.Tiploc )
  PublicTimeWrite( c, l.Pta )
  PublicTimeWrite( c, l.Ptd )
  WorkingTimeWrite( c, l.Wta )
  WorkingTimeWrite( c, l.Wtd )
  WorkingTimeWrite( c, l.Wtp )
  c.WriteString( l.Platform ).
    WriteStringArray( l.Activity ).
    WriteString( l.Line ).
    WriteString( l.Path ).
    WriteString( l.EngAllow ).
    WriteString( l.PathAllow ).
    WriteString( l.PerfAllow )
}

// BinaryCodec reader
func (l *Location) Read( c *codec.BinaryCodec ) {
  c.ReadString( &l.Id ).
    ReadString( &l.Location ).
    ReadString( &l.Tiploc )
  l.Pta = PublicTimeRead( c )
  l.Ptd = PublicTimeRead( c )
  l.Wta = WorkingTimeRead( c )
  l.Wtd = WorkingTimeRead( c )
  l.Wtp = WorkingTimeRead( c )
  c.ReadString( &l.Platform ).
    ReadStringArray( &l.Activity ).
    ReadString( &l.Line ).
    ReadString( &l.Path ).
    ReadString( &l.EngAllow ).
    ReadString( &l.PathAllow ).
    ReadString( &l.PerfAllow )
}

func (s *Schedule) appendLocation(l *Location) {
  s.Locations = append( s.Locations, l )
}
