package cif

import (
  "encoding/json"
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/rest"
  "time"
  "log"
)

// A train schedule
type Schedule struct {
  XMLName                   xml.Name  `json:"-" xml:"schedule"`
  ID struct {
    // The train UID
    TrainUID                  string    `json:"uid" xml:"uid,attr"`
    // The STP Indicator
    STPIndicator              string    `json:"stp" xml:"stp,attr"`
    // The identity sometimes confusingly called the Headcode of the service.
    // This is the value you would see in the nrod-td feed
    TrainIdentity             string    `json:"trainIdentity,omitempty" xml:"trainIdentity,attr,omitempty"`
    // The headcode of this service. Don't confuse with TrainIdentity above
    Headcode                  int       `json:"headcode,omitempty" xml:"headcode,attr,omitempty"`
  } `json:"id"`
  Runs struct {
    // The date range the schedule is valid on
    RunsFrom                  time.Time `json:"runsFrom" xml:"from,attr"`
    RunsTo                    time.Time `json:"runsTo" xml:"to,attr"`
    // The day's of the week the service will run
    DaysRun                   string    `json:"daysRun" xml:"daysRun,attr"`
    BankHolRun                string    `json:"bankHolRun,omitempty" xml:"bankHolRun,attr,omitempty"`
  } `json:"runs"`
  Meta struct {
    Status                    string    `json:"status" xml:"status,attr"`
    Category                  string    `json:"category" xml:"category,attr"`
    // The operator of this service
    ATOCCode                  string    `json:"operator,omitempty" xml:"operator,attr,omitempty"`
    ApplicableTimetable       bool      `json:"applicableTimetable" xml:"applicableTimetable,attr"`
    UICCode                   int       `json:"uic,omitempty" xml:"uic,attr,omitempty"`
    ServiceCode               int       `json:"serviceCode,omitempty" xml:"serviceCode,attr,omitempty"`
    PortionId                 string    `json:"portionId,omitempty" xml:"portionId,attr,omitempty"`
    PowerType                 string    `json:"powerType,omitempty" xml:"powerType,attr,omitempty"`
    TimingLoad                string    `json:"timingLoad,omitempty" xml:"timingLoad,attr,omitempty"`
    Speed                     int       `json:"speed,omitempty" xml:"speed,attr,omitempty"`
    OperatingCharacteristics  string    `json:",omitempty" xml:",omitempty"`
    SeatingClass              string    `json:"seatingClass,omitempty" xml:"seatingClass,attr,omitempty"`
    Sleepers                  string    `json:"sleepers,omitempty" xml:"sleepers,attr,omitempty"`
    Reservations              string    `json:"reservations,omitempty" xml:"reservations,attr,omitempty"`
    CateringCode              string    `json:"cateringCode,omitempty" xml:"cateringCode,attr,omitempty"`
    ServiceBranding           string    `json:"branding,omitempty" xml:"branding,attr,omitempty"`
  } `json:"meta"`
  // LO, LI & LT entries
  Locations              []*Location  `json:"schedule" xml:"location"`
  // The CIF extract this entry is from
  DateOfExtract             time.Time `json:"date" xml:"date,attr"`
  // URL for this Schedule
  Self                      string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

// Key returns the internal key for this schedule
func ( s *Schedule ) Key() string {
  return s.ID.TrainUID + s.Runs.RunsFrom.Format( Date ) + s.ID.STPIndicator
}

// Equals returns true if two Schedule struts refer to the same schedule.
// This checks the "primary key" for schedules which is TrainUID, RunsFrom & STPIndicator
func (s *Schedule) Equals( o *Schedule ) bool {
  if o == nil {
    return false
  }
  return s.ID.TrainUID == o.ID.TrainUID && s.Runs.RunsFrom == o.Runs.RunsFrom && s.ID.STPIndicator == o.ID.STPIndicator
}

// String returns the "primary key" for schedules which is TrainUID, RunsFrom & STPIndicator
func (s *Schedule) String() string {
  return fmt.Sprintf(
    "Schedule[uid=%s, from=%s, stp=%s]",
    s.ID.TrainUID,
    s.Runs.RunsFrom.Format( Date ),
    s.ID.STPIndicator )
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

  sj, err := json.Marshal( s )
  if err != nil {
    return err
  }

  _, err = c.tx.Exec( "SELECT timetable.addschedule( $1 )", sj )
  if err != nil {
    log.Printf( "Entry that failed:\n%s", string(sj) )
  }
  return err
}

func (c *CIF) deleteSchedule( s *Schedule ) error {
  _, err := c.tx.Exec(
    "DELETE FROM timetable.schedule WHERE uid = $1 AND stp = $2 AND startdate = $3",
    s.ID.TrainUID,
    s.ID.STPIndicator,
    s.Runs.RunsFrom.Format( Date ),
  )
  return err
}

// SetSelf sets the Schedule's Self field according to the inbound request.
// The resulting URL should then refer back to the rest endpoint that would
// return this Schedule.
func (s *Schedule) SetSelf( r *rest.Rest ) {
  s.Self = r.Self( fmt.Sprintf(
    "/schedule/%s/%s/%s",
    s.ID.TrainUID,
    s.Runs.RunsFrom.Format( Date ),
    s.ID.STPIndicator ) )
}
