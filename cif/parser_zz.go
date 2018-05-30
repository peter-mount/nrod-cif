package cif

func (c *CIF) parseZZ() error {

    // Save last schedule
    if err := c.addSchedule(); err != nil {
      return err
    }

    // Set the header to the new imported CIF
    c.header = c.importhd
    c.importhd = nil

    return nil
}
