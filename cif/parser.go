// NR CIF Parser
package cif

import (
  "errors"
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
  if err := c.parseFileHeader( scanner ); err != nil {
    return err
  }
  return c.parseFileBody( scanner )
}

// Sets the CIF structure up to the current transaction
func (c *CIF) parserInit( tx *bolt.Tx ) {
  c.tx = tx
  c.tiploc = tx.Bucket( []byte("Tiploc") )
  c.crs = tx.Bucket( []byte("Crs") )
  c.stanox = tx.Bucket( []byte("Stanox") )
  c.schedule = tx.Bucket( []byte("Schedule") )
  c.curSchedule = nil
  c.update = false
}

// Looks for the initial HD record
// If the CIF is for a full import then reset the DB
func (c *CIF) parseFileHeader( scanner *bufio.Scanner ) error {
  return c.db.Update( func( tx *bolt.Tx ) error {

    c.parserInit( tx )

    if scanner.Scan() {

      line := scanner.Text()

      if line[0:2] == "HD" {

        if err := c.parseHD( line ); err != nil {
          return err
        }

        if !c.header.Update {
          if err := c.resetDB(); err != nil {
            return err
          }
        }

        return nil
      }
    }

    return errors.New( "Not a CIF file" )
  })
}

// Parses the rest of the file after the header
func (c *CIF) parseFileBody( scanner *bufio.Scanner ) error {

  // Now run the rest of the import
  if err := c.db.Update( func( tx *bolt.Tx ) error {

    c.parserInit( tx )

    for scanner.Scan() {
      line := scanner.Text()
      if err := c.parseLine( line ); err != nil {
        return err
      }
    }
    return nil
  }); err != nil {
    return err
  }

  return nil
}

func (c *CIF) parseLine( line string ) error {
  switch line[0:2] {
    case "HD":
      return errors.New( "Erroneous HD record encountered" )

    case "TI":
      return c.parseTI( line )

    case "TA":
      return c.parseTA( line )

    case "TD":
      return c.parseTD( line )

    case "BS":
      return c.parseBS( line )

    case "BX":
      return c.parseBX( line )

    case "LO":
      return c.parseLO( line )

    case "LI":
      return c.parseLI( line )

    case "LT":
      return c.parseLT( line )

    case "ZZ":
      return c.parseZZ()

    // Ignore any unsupported records
    default:
      return nil
  }
}
