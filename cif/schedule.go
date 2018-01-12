package cif

import (
  "fmt"
  "github.com/peter-mount/golib/codec"
  "time"
)

// A train schedule
type Schedule struct {
  // The train UID
  TrainUID                  string
  // The date range the schedule is valid on
  RunsFrom                  time.Time
  RunsTo                    time.Time
  // The day's of the week the service will run
  DaysRun                   string
  BankHolRun                string    `json:",omitempty"`
  Status                    string
  Category                  string
  // The identity sometimes confusingly called the Headcode of the service.
  // This is the value you would see in the nrod-td feed
  TrainIdentity             string    `json:",omitempty"`
  // The headcode of this service. Don't confuse with TrainIdentity above
  Headcode                  int       `json:",omitempty"`
  ServiceCode               int       `json:",omitempty"`
  PortionId                 string    `json:",omitempty"`
  PowerType                 string    `json:",omitempty"`
  TimingLoad                string    `json:",omitempty"`
  Speed                     int       `json:",omitempty"`
  OperatingCharacteristics  string    `json:",omitempty"`
  SeatingClass              string    `json:",omitempty"`
  Sleepers                  string    `json:",omitempty"`
  Reservations              string    `json:",omitempty"`
  CateringCode              string    `json:",omitempty"`
  ServiceBranding           string    `json:",omitempty"`
  // The STP Indicator
  STPIndicator              string
  UICCode                   int       `json:",omitempty"`
  // The operator of this service
  ATOCCode                  string    `json:",omitempty"`
  ApplicableTimetable       bool
  // LO, LI & LT entries
  Locations              []*Location
  // The CIF extract this entry is from
  DateOfExtract             time.Time
}

// BinaryCodec writer
func ( s *Schedule) Write( c *codec.BinaryCodec ) {
  c.WriteString( s.TrainUID ).
    WriteTime( s.RunsFrom ).
    WriteTime( s.RunsTo ).
    WriteString( s.DaysRun).
    WriteString( s.BankHolRun).
    WriteString( s.Status).
    WriteString( s.Category).
    WriteString( s.TrainIdentity).
    WriteInt( s.Headcode).
    WriteInt( s.ServiceCode).
    WriteString( s.PortionId).
    WriteString( s.PowerType).
    WriteString( s.TimingLoad).
    WriteInt( s.Speed).
    WriteString( s.OperatingCharacteristics).
    WriteString( s.SeatingClass).
    WriteString( s.Sleepers).
    WriteString( s.Reservations).
    WriteString( s.CateringCode).
    WriteString( s.ServiceBranding).
    WriteString( s.STPIndicator).
    WriteInt( s.UICCode).
    WriteString( s.ATOCCode).
    WriteBool( s.ApplicableTimetable).
    WriteTime( s.DateOfExtract )

  c.WriteInt16( int16( len( s.Locations ) ) )
  for _, l := range s.Locations {
    c.Write( l )
  }
}

// BinaryCodec reader
func ( s *Schedule) Read( c *codec.BinaryCodec ) {
  c.ReadString( &s.TrainUID ).
    ReadTime( &s.RunsFrom ).
    ReadTime( &s.RunsTo ).
    ReadString( &s.DaysRun).
    ReadString( &s.BankHolRun).
    ReadString( &s.Status).
    ReadString( &s.Category).
    ReadString( &s.TrainIdentity).
    ReadInt( &s.Headcode).
    ReadInt( &s.ServiceCode).
    ReadString( &s.PortionId).
    ReadString( &s.PowerType).
    ReadString( &s.TimingLoad).
    ReadInt( &s.Speed).
    ReadString( &s.OperatingCharacteristics).
    ReadString( &s.SeatingClass).
    ReadString( &s.Sleepers).
    ReadString( &s.Reservations).
    ReadString( &s.CateringCode).
    ReadString( &s.ServiceBranding).
    ReadString( &s.STPIndicator).
    ReadInt( &s.UICCode).
    ReadString( &s.ATOCCode).
    ReadBool( &s.ApplicableTimetable).
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
  return s.TrainUID == o.TrainUID && s.RunsFrom == o.RunsFrom && s.STPIndicator == o.STPIndicator
}

// String returns the "primary key" for schedules which is TrainUID, RunsFrom & STPIndicator
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

  //
  key := []byte( s.TrainUID + s.RunsFrom.Format( Date ) + s.STPIndicator )

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

  key := []byte( s.TrainUID + s.RunsFrom.Format( Date ) + s.STPIndicator )

  return c.schedule.Delete( key )
}
