package cif

func (c *CIF ) parseBX( l string ) error {
  s := c.curSchedule

  i := 2
  i+=4 // traction class
  i = parseInt( l, i, 5, &s.UICCode )
  i = parseString( l, i, 2, &s.ATOCCode )

  var atc string
  i = parseString( l, i, 1, &atc )
  s.ApplicableTimetable = atc == "Y"

  return nil
}
