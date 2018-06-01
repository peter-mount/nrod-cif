package cifimport

func (c *CIFImporter) parseZZ() error {

    // Save last schedule
    if err := c.addSchedule(); err != nil {
      return err
    }

    // Set the header to the new imported CIF
    c.header = c.importhd
    c.importhd = nil

    return nil
}
