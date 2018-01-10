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
