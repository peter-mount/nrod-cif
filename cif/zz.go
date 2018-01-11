package cif

func (c *CIF) parseZZ() error {

    // Save last schedule
    if err := c.addSchedule(); err != nil {
      return err
    }

    // Rebuild any indices
    //if err := c.Rebuild( c.tx ); err != nil {
    //  return err
    //}

    // Finally update header to this imported one
    if err := c.put( c.tx.Bucket( []byte("Meta") ), "lastCif", c.importhd ); err != nil {
      return err
    }

    // Set the header to the new imported CIF
    c.header = c.importhd
    c.importhd = nil

    return nil
}
