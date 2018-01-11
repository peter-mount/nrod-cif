package cif

import (
  "fmt"
  "log"
  "sort"
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

func (c *CIF ) parseBS( l string ) {
  tx := l[2:3]

  switch tx {
    // New entry
    case "N":
      c.parseBSNew( l )

    // Revise - treat as new as we ensure only a single instance
    case "R":
      c.parseBSNew( l )

    // Delete
    case "D":
      c.parseBSDelete( l )
  }

}

func (c *CIF ) parseBSNew( l string ) {
  s := c.curSchedule

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
}

func (c *CIF ) parseBX( l string ) {
  s := c.curSchedule

  i := 2
  i+=4 // traction class
  i = parseInt( l, i, 5, &s.UICCode )
  i = parseString( l, i, 2, &s.ATOCCode )

  var atc string
  i = parseString( l, i, 1, &atc )
  s.ApplicableTimetable = atc == "Y"
}

func (c *CIF ) parseBSDelete( l string ) *Schedule {
  var s Schedule = Schedule{}
  i :=2
  i++ // tx
  i = parseString( l, i, 6, &s.TrainUID )
  i = parseYYMMDD( l, i, &s.RunsFrom )
  parseString( l, 79, 1, &s.STPIndicator )

  c.deleteSchedule( &s )

  return nil
}

func (c *CIF) parseLO( l string ) {
  s := c.curSchedule

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

func (c *CIF) parseLI( l string ) {
  s := c.curSchedule

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

func (c *CIF) parseLT( l string ) {
  s := c.curSchedule

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

func (c *CIF) addSchedule() error {
  s := c.curSchedule

  var ar []*Schedule

  c.get( c.schedule, s.TrainUID, &ar )

  // Check to see if we have a comparable entry. If so then replace it
  for i, e := range ar {
    if s.Equals( e ) {
      ar[ i ] = s
      return c.put( c.schedule, s.TrainUID, ar )
    }
  }

  // It's new for this uid so append & persist
  ar = append( ar, s )
  return c.put( c.schedule, s.TrainUID, ar )
}

func (c *CIF) deleteSchedule( s *Schedule ) error {
  var ar []*Schedule

  c.get( c.schedule, s.TrainUID, &ar )

  // Form a new slice without the schedule
  var n []*Schedule
  for _, e := range ar {
    if !s.Equals( e ) {
      n = append( n, e )
    }
  }

  // Persist or delete if the new slice is empty
  if len( n ) > 0 {
    return c.put( c.schedule, s.TrainUID, ar )
  } else {
    return c.schedule.Delete( []byte( s.TrainUID ) )
  }

  return nil
}

func (c *CIF) cleanupSchedules() error {
  log.Println( "Rebuilding Schedule bucket")

  return c.schedule.ForEach( func( k, v []byte ) error {
    var ar []*Schedule
    if err := getInterface( v, &ar ); err != nil {
      return err
    }

    sort.SliceStable( ar, func( i, j int ) bool {
      return ar[i].RunsFrom.Before( ar[j].RunsFrom ) && ar[i].STPIndicator < ar[i].STPIndicator
    })

    return c.put( c.schedule, ar[0].TrainUID, ar )
  } )

}

/*
// Returns all schedules for a train uid
func (c *CIF) GetSchedules( uid string ) []*Schedule {
  return c.schedules[ uid ]
}
*/
