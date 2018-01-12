package cif

import (
  "errors"
  bolt "github.com/coreos/bbolt"
  "bufio"
  "io"
  "log"
  "os"
)

// ImportFile imports a uncompressed CIF file retrieved from NetworkRail
// into the cif database.
// If this file is a full export then the database will be cleared first.
//
// The CIF.Mode field determines how this import is performed.
// This field is a bitmask so one or more options can be included.
// They are:
//
// TIPLOC   Import tiplocs
//
// SCHEDULE Import schedules
//
// ALL      Import everything, the default and the same as TIPLOC | SCHEDULE
func (c *CIF) ImportFile( fname string ) error {
  file,err := os.Open( fname )
  if err != nil {
    return err
  }

  defer file.Close()

  return c.ImportCIF( file )
}

// ImportCIF imports a uncompressed CIF file retrieved from NetworkRail
// into the cif database.
// If this file is a full export then the database will be cleared first.
//
// The CIF.Mode field determines how this import is performed.
// This field is a bitmask so one or more options can be included.
// They are:
//
// TIPLOC   Import tiplocs
//
// SCHEDULE Import schedules
//
// ALL      Import everything, the default and the same as TIPLOC | SCHEDULE
func (c *CIF) ImportCIF( r io.Reader ) error {
  scanner := bufio.NewScanner( r )

  if err := c.parseFile( scanner ); err != nil {
    return err
  }

  return scanner.Err()
}

func (c *CIF) parseFile( scanner *bufio.Scanner ) error {
  if c.Mode == 0 {
    c.Mode = TIPLOC | SCHEDULE
  }

  var lastLine string

  // Parse the header in it's own tx. This may wipe the DB if its a full import
  err := c.parseFileHeader( scanner )
  if err != nil {
    return err
  }

  if (c.Mode & TIPLOC) == TIPLOC {
    lastLine, err = c.parseTiplocs( scanner )
  }
  if err != nil {
    return err
  }

  if (c.Mode & SCHEDULE) == SCHEDULE {
    err = c.parseSchedules( scanner, lastLine )
  }

  return err
}

// Sets the CIF structure up to the current transaction
func (c *CIF) parserInit( tx *bolt.Tx ) {
  c.tx = tx

  if (c.Mode & TIPLOC) == TIPLOC {
    c.tiploc = tx.Bucket( []byte("Tiploc") )
    c.crs = tx.Bucket( []byte("Crs") )
    c.stanox = tx.Bucket( []byte("Stanox") )
  }

  if (c.Mode & SCHEDULE) == SCHEDULE {
    c.schedule = tx.Bucket( []byte("Schedule") )
    c.curSchedule = nil
  }

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
func (c *CIF) parseTiplocs( scanner *bufio.Scanner ) ( string, error ) {
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
    return "", err
  }

  // Now rebuild the Tiploc based buckets
  if err := c.db.Update( func( tx *bolt.Tx ) error {
    c.parserInit( tx )

    if err := c.cleanupStanox(); err != nil {
      return err
    }

    return c.cleanupCRS()
  }); err != nil {
    return "", err
  }

  return lastLine, nil
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

    count := 0

    c.parserInit( tx )

    // process the last line then continue & process that line
    if lastLine != "" {
      if err := c.parseSchedule( lastLine, &count ); err != nil {
        return err
      }
    }

    for scanner.Scan() {
      line := scanner.Text()
      if err := c.parseSchedule( line, &count ); err != nil {
        return err
      }
    }

    return nil
  }); err != nil {
    return err
  }

  return nil
}

func (c *CIF) parseSchedule( line string, count *int ) error {
  switch line[0:2] {
    case "BS":
      *count++
      if (*count % 25000) == 0 {
        log.Println( "Read", *count)
      }
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
      log.Println( "Read", *count)
      return c.parseZZ()

    // Ignore any unsupported records
    default:
      return nil
  }
}
