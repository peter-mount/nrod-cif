package cif

func (c *CIF ) parseBS( l string ) error {

  // Persist the last schedule as its now complete
  if err := c.addSchedule(); err != nil {
    return err
  }

  // Switch against the transaction
  switch l[2:3] {
    // New entry
    case "N":
      c.parseBSNew( l, false )

    // Revise - treat as new as we ensure only a single instance
    case "R":
      c.parseBSNew( l, true )

    // Delete
    case "D":
      c.parseBSDelete( l )
  }

  return nil
}

func (c *CIF ) parseBSNew( l string, update bool ) {
  s := &Schedule{}
  c.curSchedule = s
  c.update = update

  // Skip BS
  i := 2
  i++ // TX
  i = parseString( l, i, 6, &s.TrainUID )
  i = parseYYMMDD( l, i, &s.RunsFrom )
  i = parseYYMMDD( l, i, &s.RunsTo )
  i = parseString( l, i, 7, &s.DaysRun )
  i = parseStringTrim( l, i, 1, &s.BankHolRun )
  i = parseString( l, i, 1, &s.Status )
  i = parseString( l, i, 2, &s.Category )
  i = parseStringTrim( l, i, 4, &s.TrainIdentity )
  i = parseInt( l, i, 4, &s.Headcode )
  i++ // Course Indicator
  i = parseInt( l, i, 8, &s.ServiceCode )
  i = parseStringTrim( l, i, 1, &s.PortionId )
  i = parseStringTrim( l, i, 3, &s.PowerType )
  i = parseStringTrim( l, i, 4, &s.TimingLoad )
  i = parseInt( l, i, 3, &s.Speed )
  i = parseStringTrim( l, i, 6, &s.OperatingCharacteristics )
  i = parseStringTrim( l, i, 1, &s.SeatingClass )
  i = parseStringTrim( l, i, 1, &s.Sleepers )
  i = parseStringTrim( l, i, 1, &s.Reservations )
  i++ // Connection Indicator
  i = parseStringTrim( l, i, 4, &s.CateringCode )
  i = parseStringTrim( l, i, 4, &s.ServiceBranding )
  i++ // Spare
  i = parseString( l, i, 1, &s.STPIndicator )
}

func (c *CIF ) parseBSDelete( l string ) *Schedule {
  c.curSchedule = nil

  var s Schedule = Schedule{}

  i :=2
  i++ // tx
  i = parseString( l, i, 6, &s.TrainUID )
  i = parseYYMMDD( l, i, &s.RunsFrom )
  parseString( l, 79, 1, &s.STPIndicator )

  c.deleteSchedule( &s )

  return nil
}
