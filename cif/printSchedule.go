// Debug function to return a schedule as a String
// Unlike String() this returns the full schedule
package cif

import (
  "github.com/peter-mount/golib/util"
  "strings"
)

func (s *Schedule) FullString() string {
  var b util.BufferWriter

  b.Field( "UID", s.TrainUID )
  b.Field( "Runs From", s.RunsFrom.Format( HumanDate ) )
  b.Field( "To", s.RunsTo.Format( HumanDate ) )
  b.Field( "STPIndicator", s.STPIndicator )
  b.Field( "Bank Hol Run", s.BankHolRun )
  b.Field( "Status", s.Status )
  b.Field( "Category", s.Category )
  b.Field( "Identity", s.TrainIdentity )
  b.FieldInt( "Headcode", s.Headcode )
  b.FieldInt( "UIC Code", s.UICCode )
  b.Field( "ATOCCode", s.ATOCCode )
  b.FieldBool( "ApplicableTimetable", s.ApplicableTimetable )
  b.FieldInt( "Service Code", s.ServiceCode )
  b.Field( "PortionId", s.PortionId )
  b.Field( "PowerType", s.PowerType )
  b.Field( "TimingLoad", s.TimingLoad )
  b.FieldInt( "Speed", s.Speed )
  b.Field( "OperatingCharacteristics", s.OperatingCharacteristics )
  b.Field( "SeatingClass", s.SeatingClass )
  b.Field( "Sleepers", s.Sleepers )
  b.Field( "Reservations", s.Reservations )
  b.Field( "CateringCode", s.CateringCode )
  b.Field( "ServiceBranding", s.ServiceBranding )

  b.Row()
  b.Pad( "", 2 )
  b.Pad( "ID", 2)
  b.Pad( "Location", 8 )
  b.Pad( "PTA", 5 ).Pad( "PTD", 5 )
  b.Pad( "WTA", 8 ).Pad( "WTD", 8 ).Pad( "WTP", 8 )
  b.Pad( "Plat", 4 )
  b.Pad( "Line", 4 )
  b.Pad( "Path", 4 )
  b.Pad( "Eng", 4 ).Pad( "Path", 4 ).Pad( "Perf", 4 )
  b.Pad( "Activity", 12 )

  for i, l := range s.Locations {
    b.Row()
    b.PadInt( i, 2 )
    b.Pad( l.Id, 2 )
    b.Pad( l.Location, 8 )
    b.PadX( 5 ).PadX( 5 )
    b.PadX( 8 ).PadX( 8 ).PadX( 8 )
    b.Pad( l.Platform, 4 )
    b.Pad( l.Line, 4 )
    b.Pad( l.Path, 4 )
    b.Pad( l.EngAllow, 4 ).Pad( l.PathAllow, 4 ).Pad( l.PerfAllow, 4 )
    b.Pad( strings.Join( l.Activity, "" ), 12 )
  }

  return b.String()
}
