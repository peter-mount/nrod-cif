// NR CIF Parser
package cif

import (
  "errors"
  "bufio"
  bolt "github.com/coreos/bbolt"
  "log"
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
  // Parse the header in it's own tx. This may wipe the DB if its a full import
  if err := c.parseFileHeader( scanner ); err != nil {
    return err
  }

  // Now start with Tiplocs
  return c.parseTiplocs( scanner )
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
func (c *CIF) parseTiplocs( scanner *bufio.Scanner ) error {
  log.Println( "Parsing Tiploc's" )

  var lastLine string

  // Now run the rest of the import
  if err := c.db.Update( func( tx *bolt.Tx ) error {

    c.parserInit( tx )

    for scanner.Scan() {
      line := scanner.Text()
      if bail, err := c.parseTiploc( line ); err != nil {
        return err
      } else if bail {
        lastLine = line
        return nil
      }
    }

    return nil
  }); err != nil {
    return err
  }

  // Now rebuild the Tiploc based buckets
  if err := c.db.Update( func( tx *bolt.Tx ) error {
    c.parserInit( tx )

    if err := c.cleanupStanox(); err != nil {
      return err
    }

    return c.cleanupCRS()
  }); err != nil {
    return err
  }

  // Procede to the next block
  return c.parseSchedules( scanner, lastLine )
}

func (c *CIF) parseTiploc( line string ) ( bool, error ) {
  switch line[0:2] {
    case "TI":
      return false, c.parseTI( line )

    case "TA":
      return false, c.parseTA( line )

    case "TD":
      return false, c.parseTD( line )

    // If not a Tiploc record then bail out to the next stage
    default:
      return true, nil
  }
}

// Parses the rest of the file after the header
func (c *CIF) parseSchedules( scanner *bufio.Scanner, lastLine string ) error {
  log.Println( "Parsing Schedules" )

  // Now run the rest of the import
  if err := c.db.Update( func( tx *bolt.Tx ) error {

    c.parserInit( tx )

    // process the last line then continue & process that line
    if lastLine != "" {
      if err := c.parseSchedule( lastLine ); err != nil {
        return err
      }
    }

    for scanner.Scan() {
      line := scanner.Text()
      if err := c.parseSchedule( line ); err != nil {
        return err
      }
    }

    return nil
  }); err != nil {
    return err
  }

  return nil
}

func (c *CIF) parseSchedule( line string ) error {
  switch line[0:2] {
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
