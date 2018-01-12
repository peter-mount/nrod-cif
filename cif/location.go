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
// fields. Tiploc is the name of this location.
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
  Id          string
  // Location including Suffix (for circular routes)
  // This is guaranteed to be unique per schedule, although for most purposes
  // like display you would use Tiploc
  Location    string
  // Tiploc of this location. For some schedules like circular routes this can
  // appear more than once in a schedule.
  Tiploc      string
  // Public Timetable
  Pta         PublicTime
  Ptd         PublicTime
  // Working Timetable
  Wta         WorkingTime
  Wtd         WorkingTime
  Wtp         WorkingTime
  // Platform
  Platform    string
  // Activity up to 6 codes
  Activity  []string
  // The Line the train will take
  Line        string
  // The Path the train will take
  Path        string
  // Allowances at this location
  EngAllow    string
  PathAllow   string
  PerfAllow   string
}

func (l *Location) Write( c *codec.BinaryCodec ) {
  c.WriteString( l.Id ).
    WriteString( l.Location ).
    WriteString( l.Tiploc ).
    Write( &l.Pta ).
    Write( &l.Ptd ).
    Write( &l.Wta ).
    Write( &l.Wtd ).
    Write( &l.Wtp ).
    WriteString( l.Platform ).
    WriteStringArray( l.Activity ).
    WriteString( l.Line ).
    WriteString( l.Path ).
    WriteString( l.EngAllow ).
    WriteString( l.PathAllow ).
    WriteString( l.PerfAllow )
}

func (l *Location) Read( c *codec.BinaryCodec ) {
  c.ReadString( &l.Id ).
    ReadString( &l.Location ).
    ReadString( &l.Tiploc ).
    Read( &l.Pta ).
    Read( &l.Ptd ).
    Read( &l.Wta ).
    Read( &l.Wtd ).
    Read( &l.Wtp ).
    ReadString( &l.Platform ).
    ReadStringArray( &l.Activity ).
    ReadString( &l.Line ).
    ReadString( &l.Path ).
    ReadString( &l.EngAllow ).
    ReadString( &l.PathAllow ).
    ReadString( &l.PerfAllow )
}

func newLocation() *Location {
  var loc *Location = &Location{}
  loc.Pta.Set( -1 )
  loc.Ptd.Set( -1 )
  loc.Wta.Set( -1 )
  loc.Wtd.Set( -1 )
  loc.Wtp.Set( -1 )
  return loc
}

func (s *Schedule) appendLocation(l *Location) {
  s.Locations = append( s.Locations, l )
}
