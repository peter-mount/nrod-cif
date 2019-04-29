package cifimport

func (c *CIFImporter) parseBX(l string) error {
  s := c.curSchedule

  i := 2
  i += 4 // traction class
  i = parseInt(l, i, 5, &s.Meta.UICCode)
  i = parseStringTrim(l, i, 2, &s.Meta.ATOCCode)

  var atc string
  i = parseString(l, i, 1, &atc)
  s.Meta.ApplicableTimetable = atc == "Y"

  return nil
}
