package cif

import (
  "log"
  "fmt"
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

}

func (s *Schedule) String() string {
  return fmt.Sprintf(
    "Schedule %s %s %s",
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
      return c.parseBSRevise( l )

    // Delete
    case "D":
      return c.parseBSDelete( l )
  }

  return nil
}

func (c *CIF ) parseBSNew( l string ) *Schedule {
  var s *Schedule = &Schedule{}

  // Skip BSN
  i := 3
  i = parseString( l, i, 6, &s.TrainUID )
  i = parseYYMMDD( l, i, &s.RunsFrom )
  i = parseYYMMDD( l, i, &s.RunsTo )
  i = parseString( l, i, 7, &s.DaysRun )
  i = parseString( l, i, 1, &s.BankHolRun )
  i = parseString( l, i, 1, &s.Status )
  i = parseString( l, i, 2, &s.Category )
  i = parseString( l, i, 4, &s.TrainIdentity )
  i = parseInt( l, i, 1, &s.Headcode )
  i++ // Course Indicator
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

func (c *CIF ) parseBSRevise( l string ) *Schedule {
  s := c.parseBSNew( l )
  // todo delete existing entry?
  return s
}

func (c *CIF ) parseBSDelete( l string ) *Schedule {
  log.Fatal( "Delete not yet implemented" )
  return nil
}

// Returns all schedules for a train uid
func (c *CIF) GetSchedules( uid string ) []*Schedule {
  return c.schedules[ uid ]
}

func (c *CIF) addSchedule( s *Schedule ) {
  //ary := c.schedules[ s.TrainUID ]

}
