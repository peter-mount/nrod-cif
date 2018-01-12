package cif

import (
  "github.com/peter-mount/golib/codec"
)

type Location struct {
  // LO,
  Id          string
  // Location including Suffix (for circular routes)
  Location    string
  // Tiploc of Location sans Suffix
  Tiploc      string
  // Public times in seconds of day
  Pta         PublicTime
  Ptd         PublicTime
  // Working times in seconds of day
  Wta         WorkingTime
  Wtd         WorkingTime
  Wtp         WorkingTime
  // Platform
  Platform    string
  // Activity up to 6 codes
  Activity  []string
  // Misc
  Line        string
  Path        string
  EngAllow    string
  PathAllow   string
  PerfAllow   string
}

func (l *Location) Write( c *codec.BinaryCodec ) {
  c.WriteString( l.Id ).
    WriteString( l.Location ).
    WriteString( l.Tiploc ).
    Write( l.Pta ).
    Write( l.Ptd ).
    Write( l.Wta ).
    Write( l.Wtd ).
    Write( l.Wtp ).
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
