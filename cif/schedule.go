package cif

import (
  "fmt"
  "log"
  "sort"
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
  // LO, LI & LT entries
  Locations              []*Location
  // The CIF extract this entry is from
  DateOfExtract             time.Time
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

func (c *CIF) addSchedule() error {
  // Do nothing if we have no schedule to persist
  if c.curSchedule == nil {
    return nil
  }

  // get schedule & reset
  s := c.curSchedule
  c.curSchedule = nil

  // Link it to this CIF file & persist
  s.DateOfExtract = c.importhd.DateOfExtract

  var ar []*Schedule

  c.get( c.schedule, s.TrainUID, &ar )

  // Check to see if we have a comparable entry. If so then replace it
  for i, e := range ar {
    if s.Equals( e ) {
      // Only replace & persist if the new schedule was extracted after the existing entry
      if s.DateOfExtract.After( e.DateOfExtract ) {
        ar[ i ] = s
        return c.put( c.schedule, s.TrainUID, ar )
      }
      return nil
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
