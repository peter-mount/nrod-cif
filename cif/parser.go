// NR CIF Parser
package cif

import (
  "bufio"
  bolt "github.com/coreos/bbolt"
  "os"
)

func (c *CIF) Parse( fname string ) error {
  file,err := os.Open( fname )
  if err != nil {
    return err
  }

  defer file.Close()

  scanner := bufio.NewScanner( file )
  if err := c.parseFile( scanner ); err != nil {
    return err
  }

  if err := scanner.Err(); err != nil {
    return err
  }

  return nil
}

func (c *CIF) parseFile( scanner *bufio.Scanner ) error {
  if err := c.db.Update( func( tx *bolt.Tx ) error {

    c.tx = tx
    c.tiploc = tx.Bucket( []byte("Tiploc") )
    c.crs = tx.Bucket( []byte("Crs") )
    c.stanox = tx.Bucket( []byte("Stanox") )
    c.schedule = tx.Bucket( []byte("Schedule") )

    var schedule *Schedule

    for scanner.Scan() {
      line := scanner.Text()
      switch line[0:2] {
      case "HD":
        if err := c.parseHD( line ); err != nil {
          return err
        }

        if !c.header.Update {
          if err := c.resetDB(); err != nil {
            return err
          }
        }

      case "TI":
        if err := c.parseTiplocInsert( line ); err != nil {
          return err
        }

      case "TA":
        if err := c.parseTiplocAmend( line ); err != nil {
          return err
        }

      case "TD":
        if err := c.parseTiplocDelete( line ); err != nil {
          return err
        }

      case "BS":
        // Persist the last schedule as its now complete
        if c.curSchedule != nil {
          if err := c.addSchedule(); err != nil {
            return err
          }
        }

        c.curSchedule = &Schedule{}
        c.parseBS( line )

      case "BX":
        c.parseBX( line )

        /*
      case "LO":
      c.parseLO( line )

    case "LI":
    c.parseLI( line )

  case "LT":
  c.parseLT( line )
  */

      case "ZZ":
        // Save last schedule
        if schedule != nil {
          if err := c.addSchedule(); err != nil {
            return err
          }
        }

        if err := c.Rebuild( c.tx ); err != nil {
          return err
        }

        // Finally update header
        if h, err := c.GetHD(); err != nil {
          return err
        } else {
          c.header = h
        }
      }
    }
    return nil
  }); err != nil {
    return err
  }

  return nil
}
