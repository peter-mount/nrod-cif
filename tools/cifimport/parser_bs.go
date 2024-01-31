package cifimport

import (
	"github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFImporter) parseBS(l string) error {

	// Persist the last schedule as its now complete
	if err := c.addSchedule(); err != nil {
		return err
	}

	// Switch against the transaction
	switch l[2:3] {
	// New entry
	case "N":
		c.parseBSNew(l, false)

	// Revise - treat as new as we ensure only a single instance
	case "R":
		c.parseBSNew(l, true)

	// Delete
	case "D":
		c.parseBSDelete(l)
	}

	return nil
}

func (c *CIFImporter) parseBSNew(l string, update bool) {
	s := &cif.Schedule{}
	c.curSchedule = s
	c.update = update

	// Skip BS
	i := 2
	i++ // TX
	i = parseString(l, i, 6, &s.ID.TrainUID)
	i = parseYYMMDD(l, i, &s.Runs.RunsFrom)
	i = parseYYMMDD(l, i, &s.Runs.RunsTo)
	i = parseString(l, i, 7, &s.Runs.DaysRun)
	i = parseStringTrim(l, i, 1, &s.Runs.BankHolRun)
	i = parseString(l, i, 1, &s.Meta.Status)
	i = parseString(l, i, 2, &s.Meta.Category)
	i = parseStringTrim(l, i, 4, &s.ID.TrainIdentity)
	i = parseInt(l, i, 4, &s.ID.Headcode)
	i++ // Course Indicator
	i = parseInt(l, i, 8, &s.Meta.ServiceCode)
	i = parseStringTrim(l, i, 1, &s.Meta.PortionId)
	i = parseStringTrim(l, i, 3, &s.Meta.PowerType)
	i = parseStringTrim(l, i, 4, &s.Meta.TimingLoad)
	i = parseInt(l, i, 3, &s.Meta.Speed)
	i = parseStringTrim(l, i, 6, &s.Meta.OperatingCharacteristics)
	i = parseStringTrim(l, i, 1, &s.Meta.SeatingClass)
	i = parseStringTrim(l, i, 1, &s.Meta.Sleepers)
	i = parseStringTrim(l, i, 1, &s.Meta.Reservations)
	i++ // Connection Indicator
	i = parseStringTrim(l, i, 4, &s.Meta.CateringCode)
	i = parseStringTrim(l, i, 4, &s.Meta.ServiceBranding)
	i++ // Spare
	i = parseString(l, i, 1, &s.ID.STPIndicator)
}

func (c *CIFImporter) parseBSDelete(l string) *cif.Schedule {
	c.curSchedule = nil

	var s cif.Schedule = cif.Schedule{}

	i := 2
	i++ // tx
	i = parseString(l, i, 6, &s.ID.TrainUID)
	i = parseYYMMDD(l, i, &s.Runs.RunsFrom)
	parseString(l, 79, 1, &s.ID.STPIndicator)

	c.deleteSchedule(&s)

	return nil
}
