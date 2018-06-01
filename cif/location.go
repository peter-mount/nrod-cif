package cif

import (
  "github.com/peter-mount/nre-feeds/util"
  "strings"
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
  Id          string        `json:"-" xml:"-"`
  // Location including Suffix (for circular routes)
  // This is guaranteed to be unique per schedule, although for most purposes
  // like display you would use Tiploc
  Location    string        `json:"-" xml:"-"`
  // Tiploc of this location. For some schedules like circular routes this can
  // appear more than once in a schedule.
  Tiploc      string        `json:"tpl" xml:"tpl,attr"`
  Visit       int           `json:"visit,omitempty" xml:"visit,attr,omitempty"`
  // Public Timetable
  Times       util.CircularTimes `json:"time"`
  // Platform
  Platform    string        `json:"plat,omitempty" xml:"plat,attr,omitempty"`
  // Activity up to 6 codes
  Activity  []string        `json:"activity,omitempty" xml:"activity,omitempty"`
  // The Line the train will take
  Line        string        `json:"line,omitempty" xml:"line,attr,omitempty"`
  // The Path the train will take
  Path        string        `json:"path,omitempty" xml:"path,attr,omitempty"`
  // Allowances at this location
  EngAllow    string        `json:"engAllow,omitempty" xml:"engAllow,attr,omitempty"`
  PathAllow   string        `json:"pathAllow,omitempty" xml:"pathAllow,attr,omitempty"`
  PerfAllow   string        `json:"perfAllow,omitempty" xml:"perfAllow,attr,omitempty"`
}

func (s *Schedule) AppendLocation(l *Location) {
  // Used in parsing, set tiploc & Visit accordingly
  if l.Tiploc == "" && l.Location != "" {
    l.Tiploc = strings.Trim( l.Location[0:7], " " )
    if len( l.Location ) > 7 && l.Location[7] != ' ' {
      l.Visit = int( l.Location[7] - '0' )
    }
    l.Location = strings.Trim( l.Location, " " )
  }

  if l.Times.Pta != nil && (*l.Times.Pta).IsZero() {
    l.Times.Pta = nil
  }

  if l.Times.Ptd != nil && (*l.Times.Ptd).IsZero() {
    l.Times.Ptd = nil
  }

  l.Times.UpdateTime()

  s.Locations = append( s.Locations, l )
}
