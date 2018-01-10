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
  //var schedule *Schedule

  for scanner.Scan() {
    line := scanner.Text()
    switch line[0:2] {
      case "HD":
        if !c.parseHD( line ) {
          // Skip this file
          return nil
        }

      case "TI":
        c.parseTiplocInsert( line )

      case "TA":
        c.parseTiplocAmend( line )

      case "TD":
        c.parseTiplocDelete( line )

      case "BS":
        //schedule =
        c.parseBS( line )

      case "ZZ":
        c.cleanup()
    }
  }

  return nil
}
