// NR CIF Parser
package cif

import (
  "bufio"
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
  tx, err := c.db.Begin(true)
  if err != nil {
    return err
  }
  defer tx.Rollback()

  c.tx = tx
  c.tiploc = tx.Bucket( []byte("Tiploc") )
  c.crs = tx.Bucket( []byte("Crs") )
  c.stanox = tx.Bucket( []byte("Stanox") )

  //var schedule *Schedule

  for scanner.Scan() {
    line := scanner.Text()
    switch line[0:2] {
      case "HD":
        if err := c.parseHD( line ); err != nil {
          return err
        }

        if !c.Header.Update {
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

        /*
      case "BS":
        schedule = c.parseBS( line )
        if schedule != nil {
          c.addSchedule( schedule )
        }

      case "BX":
        if schedule != nil {
          c.parseBX( line, schedule )
        }

      case "LO":
        if schedule != nil {
          c.parseLO( line, schedule )
        }

      case "LI":
        if schedule != nil {
          c.parseLI( line, schedule )
        }

      case "LT":
        if schedule != nil {
          c.parseLT( line, schedule )
        }
        */

      case "ZZ":
        if err := c.Rebuild( c.tx ); err != nil {
          return err
        }
    }
  }

  return tx.Commit()
}
