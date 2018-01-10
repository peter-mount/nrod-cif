package cif

import (
  "log"
  "fmt"
  "strings"
  "time"
)

type Schedule struct {
  // BS record
  TrainUID                  string
  RunsFrom                  time.Time
  RunsTo                    time.Time
  DaysRun                   string
  BankHolRun                string
  Status                    string
  Category                  string
  TrainIdentity             string
  Headcode                  int
  ServiceCode               int
  PortionId                 string
  PowerType                 string
  TimingLoad                string
  Speed                     int
  OperatingCharacteristics  string
  SeatingClass              string
  Sleepers                  string
  Reservations              string
  CateringCode              string
  ServiceBranding           string
  STPIndicator              string
  // BX record
  UICCode                   int
  ATOCCode                  string
  ApplicableTimetable       bool
  Locations              []*Location
}

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
func (s *Schedule) Equals( o *Schedule ) bool {
  if o == nil {
    return false
  }
  return s.TrainUID == o.TrainUID && s.RunsFrom == o.RunsFrom && s.STPIndicator == o.STPIndicator
}

func (s *Schedule) String() string {
  return fmt.Sprintf(
    "Schedule[uid=%s, from=%s, stp=%s]",
    s.TrainUID,
    s.RunsFrom.Format( Date ),
    s.STPIndicator )
}

func (c *CIF ) parseBS( l string ) *Schedule {
  tx := l[2:3]

  switch tx {
    // New entry
    case "N":
      return c.parseBSNew( l )

    // Revise - treat as new as we ensure only a single instance
    case "R":
      return c.parseBSNew( l )

    // Delete
    case "D":
      return c.parseBSDelete( l )
  }

  return nil
}

func (c *CIF ) parseBSNew( l string ) *Schedule {
  var s *Schedule = &Schedule{}

  // Skip BS
  i := 2
  i++ // TX
  i = parseString( l, i, 6, &s.TrainUID )
  i = parseYYMMDD( l, i, &s.RunsFrom )
  i = parseYYMMDD( l, i, &s.RunsTo )
  i = parseString( l, i, 7, &s.DaysRun )
  i = parseString( l, i, 1, &s.BankHolRun )
  i = parseString( l, i, 1, &s.Status )
  i = parseString( l, i, 2, &s.Category )
  i = parseString( l, i, 4, &s.TrainIdentity )
  i = parseInt( l, i, 4, &s.Headcode )
  i++ // Course Indicator
  i = parseInt( l, i, 8, &s.ServiceCode )
  i = parseString( l, i, 1, &s.PortionId )
  i = parseString( l, i, 3, &s.PowerType )
  i = parseString( l, i, 4, &s.TimingLoad )
  i = parseInt( l, i, 3, &s.Speed )
  i = parseString( l, i, 6, &s.OperatingCharacteristics )
  i = parseString( l, i, 1, &s.SeatingClass )
  i = parseString( l, i, 1, &s.Sleepers )
  i = parseString( l, i, 1, &s.Reservations )
  i++ // Connection Indicator
  i = parseString( l, i, 4, &s.CateringCode )
  i = parseString( l, i, 4, &s.ServiceBranding )
  i++ // Spare
  i = parseString( l, i, 1, &s.STPIndicator )

  return s
}

func (c *CIF ) parseBX( l string, s *Schedule ) {
  i := 2
  i+=4 // traction class
  i = parseInt( l, i, 5, &s.UICCode )
  i = parseString( l, i, 2, &s.ATOCCode )

  var atc string
  i = parseString( l, i, 1, &atc )
  s.ApplicableTimetable = atc == "Y"
}

func (c *CIF ) parseBSDelete( l string ) *Schedule {
  log.Fatal( "Delete not yet implemented" )
  return nil
}

func (c *CIF) parseLO( l string, s *Schedule ) {
  var loc *Location = newLocation()
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )
  loc.Tiploc = strings.Trim( loc.Location[0:8], " " )

  i = parseHHMMS( l, i, &loc.Wtd )
  i = parseHHMM( l, i, &loc.Ptd )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Line )

  i = parseStringTrim( l, i, 2, &loc.EngAllow )
  i = parseStringTrim( l, i, 2, &loc.PathAllow )

  i = parseActivity( l, i, &loc.Activity)

  i = parseStringTrim( l, i, 2, &loc.PerfAllow )

  s.appendLocation( loc )
}

func (c *CIF) parseLI( l string, s *Schedule ) {
  var loc *Location = newLocation()
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )
  loc.Tiploc = strings.Trim( loc.Location[0:8], " " )

  i = parseHHMMS( l, i, &loc.Wta )
  i = parseHHMMS( l, i, &loc.Wtd )
  i = parseHHMMS( l, i, &loc.Wtp )

  i = parseHHMM( l, i, &loc.Pta )
  i = parseHHMM( l, i, &loc.Ptd )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Line )
  i = parseStringTrim( l, i, 3, &loc.Path )
  i = parseActivity( l, i, &loc.Activity)

  i = parseStringTrim( l, i, 2, &loc.EngAllow )
  i = parseStringTrim( l, i, 2, &loc.PathAllow )
  i = parseStringTrim( l, i, 2, &loc.PerfAllow )

  s.appendLocation( loc )
}

func (c *CIF) parseLT( l string, s *Schedule ) {
  var loc *Location = newLocation()
  i := 0
  i = parseString( l, i, 2, &loc.Id )

  // Location is Tiploc + Suffix
  i = parseString( l, i, 8, &loc.Location )
  loc.Tiploc = strings.Trim( loc.Location[0:8], " " )

  i = parseHHMMS( l, i, &loc.Wta )

  i = parseHHMM( l, i, &loc.Pta )

  i = parseStringTrim( l, i, 3, &loc.Platform )
  i = parseStringTrim( l, i, 3, &loc.Path )
  i = parseActivity( l, i, &loc.Activity )

  s.appendLocation( loc )
}

func newLocation() *Location {
  var loc *Location = &Location{}
  loc.Wta.t = -1
  loc.Wtd.t = -1
  loc.Wtp.t = -1
  return loc
}

func (s *Schedule) appendLocation(l *Location) {
  if l.Pta.t == 0 {
    l.Pta.t = -1
  }
  if l.Ptd.t == 0 {
    l.Ptd.t = -1
  }
  s.Locations = append( s.Locations, l )
}
