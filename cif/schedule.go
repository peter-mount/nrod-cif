package cif

import (
  "encoding/xml"
  "fmt"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "time"
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

// BinaryCodec writer
func ( s *Schedule) Write( c *codec.BinaryCodec ) {
  c.WriteString( s.ID.TrainUID ).
    WriteTime( s.Runs.RunsFrom ).
    WriteTime( s.Runs.RunsTo ).
    WriteString( s.Runs.DaysRun).
    WriteString( s.Runs.BankHolRun).
    WriteString( s.Meta.Status).
    WriteString( s.Meta.Category).
    WriteString( s.ID.TrainIdentity).
    WriteInt( s.ID.Headcode).
    WriteInt( s.Meta.ServiceCode).
    WriteString( s.Meta.PortionId).
    WriteString( s.Meta.PowerType).
    WriteString( s.Meta.TimingLoad).
    WriteInt( s.Meta.Speed).
    WriteString( s.Meta.OperatingCharacteristics).
    WriteString( s.Meta.SeatingClass).
    WriteString( s.Meta.Sleepers).
    WriteString( s.Meta.Reservations).
    WriteString( s.Meta.CateringCode).
    WriteString( s.Meta.ServiceBranding).
    WriteString( s.ID.STPIndicator).
    WriteInt( s.Meta.UICCode).
    WriteString( s.Meta.ATOCCode).
    WriteBool( s.Meta.ApplicableTimetable).
    WriteTime( s.DateOfExtract )

  c.WriteInt16( int16( len( s.Locations ) ) )
  for _, l := range s.Locations {
    c.Write( l )
  }
}

// BinaryCodec reader
func ( s *Schedule) Read( c *codec.BinaryCodec ) {
  c.ReadString( &s.ID.TrainUID ).
    ReadTime( &s.Runs.RunsFrom ).
    ReadTime( &s.Runs.RunsTo ).
    ReadString( &s.Runs.DaysRun).
    ReadString( &s.Runs.BankHolRun).
    ReadString( &s.Meta.Status).
    ReadString( &s.Meta.Category).
    ReadString( &s.ID.TrainIdentity).
    ReadInt( &s.ID.Headcode).
    ReadInt( &s.Meta.ServiceCode).
    ReadString( &s.Meta.PortionId).
    ReadString( &s.Meta.PowerType).
    ReadString( &s.Meta.TimingLoad).
    ReadInt( &s.Meta.Speed).
    ReadString( &s.Meta.OperatingCharacteristics).
    ReadString( &s.Meta.SeatingClass).
    ReadString( &s.Meta.Sleepers).
    ReadString( &s.Meta.Reservations).
    ReadString( &s.Meta.CateringCode).
    ReadString( &s.Meta.ServiceBranding).
    ReadString( &s.ID.STPIndicator).
    ReadInt( &s.Meta.UICCode).
    ReadString( &s.Meta.ATOCCode).
    ReadBool( &s.Meta.ApplicableTimetable).
    ReadTime( &s.DateOfExtract )

  var l int16
  c.ReadInt16( &l )
  for i := 0; i < int(l); i++ {
    loc := &Location{}
    c.Read( loc )
    s.Locations = append( s.Locations, loc )
  }
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

  //
  key := []byte( s.Key() )

  var os Schedule = Schedule{}
  b := c.schedule.Get( key )
  dec := codec.NewBinaryCodecFrom( b )
  dec.Read( &os )
  if !s.Equals( &os ) || s.DateOfExtract.After( os.DateOfExtract ) {
    enc := codec.NewBinaryCodec()
    enc.Write( s )
    if enc.Error() != nil {
      return enc.Error()
    }
    return c.schedule.Put( key, enc.Bytes() )
  }

  return nil
}

func (c *CIF) deleteSchedule( s *Schedule ) error {

  key := []byte( s.Key() )

  return c.schedule.Delete( key )
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
