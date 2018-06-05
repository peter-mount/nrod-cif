package cif

import (
  "encoding/xml"
  "time"
)

// A train schedule
type Schedule struct {
  XMLName                     xml.Name  `json:"-" xml:"schedule"`
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

// Equals returns true if two Schedule struts refer to the same schedule.
// This checks the "primary key" for schedules which is TrainUID, RunsFrom & STPIndicator
func (s *Schedule) Equals( o *Schedule ) bool {
  if o == nil {
    return false
  }
  return s.ID.TrainUID == o.ID.TrainUID && s.Runs.RunsFrom == o.Runs.RunsFrom && s.ID.STPIndicator == o.ID.STPIndicator
}
