package cif

import (
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
  // LO, LI & LT entries
  Locations              []*Location
  // The CIF extract this entry is from
  DateOfExtract             time.Time
}

func ( s *Schedule) Write( c *BinaryCodec ) {
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

func ( s *Schedule) Read( c *BinaryCodec ) {
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

  ar, _ := c.getSchedules( s.TrainUID )

  // Check to see if we have a comparable entry. If so then replace it
  for i, e := range ar {
    if s.Equals( e ) {
      // Only replace & persist if the new schedule was extracted after the existing entry
      if s.DateOfExtract.After( e.DateOfExtract ) {
        ar[ i ] = s
        return c.putSchedules( s.TrainUID, ar )
      }
      return nil
    }
  }

  // It's new for this uid so append & persist
  ar = append( ar, s )
  return c.putSchedules( s.TrainUID, ar )
}

func (c *CIF) getSchedules( uid string ) ( []*Schedule, error ) {
  var ar []*Schedule

  // Retrieve the existing entry (if any)
  b := c.schedule.Get( []byte( uid ) )
  if b != nil {
    codec := NewBinaryCodecFrom( b )

    var l int16
    codec.ReadInt16( &l )
    for i := 0; i < int(l); i++ {
      var sched *Schedule = &Schedule{}
      codec.Read( sched )
      ar = append( ar, sched )
    }

    if codec.Error() != nil {
      return nil, codec.Error()
    }
  }

  return ar, nil
}

func (c *CIF) putSchedules( uid string, ar []*Schedule ) error {
  codec := NewBinaryCodec()

  codec.WriteInt32( int32( len( ar ) ) )
  for _, s := range ar {
    codec.Write( s )
  }

  if codec.Error() != nil {
    return codec.Error()
  }

  return c.schedule.Put( []byte( uid ), codec.Bytes() )
}

func (c *CIF) deleteSchedule( s *Schedule ) error {

  ar, _ := c.getSchedules( s.TrainUID )

  // Form a new slice without the schedule
  var n []*Schedule
  for _, e := range ar {
    if !s.Equals( e ) {
      n = append( n, e )
    }
  }

  // Persist or delete if the new slice is empty
  if len( n ) > 0 {
    return c.putSchedules( s.TrainUID, ar )
  } else {
    return c.schedule.Delete( []byte( s.TrainUID ) )
  }

  return nil
}
