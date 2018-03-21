package cif

import (
  "github.com/peter-mount/golib/codec"
)

func (c *CIF) parseZZ() error {

    // Save last schedule
    if err := c.addSchedule(); err != nil {
      return err
    }

    // Finally update header to this imported one
    codec := codec.NewBinaryCodec()
    codec.Write( c.importhd )
    if codec.Error() != nil {
      return codec.Error()
    }

    if err := c.tx.Bucket( []byte("Meta") ).Put( []byte( "lastCif" ), codec.Bytes() ); err != nil {
      return err
    }

    // Set the header to the new imported CIF
    c.header = c.importhd
    c.importhd = nil

    return nil
}
