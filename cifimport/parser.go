package cifimport

import (
  "database/sql"
  "errors"
  //bolt "github.com/coreos/bbolt"
  "bufio"
  "io"
  "log"
  "os"
)

// ImportFile imports a uncompressed CIF file retrieved from NetworkRail
// into the cif database.
func (c *CIFImporter) importFile( fname string ) (bool, error) {
  file,err := os.Open( fname )
  if err != nil {
    return false, err
  }

  defer file.Close()

  skip, err := c.importCIF( file )
  return skip, err
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
func (c *CIFImporter) importCIF( r io.Reader ) (bool, error) {
  scanner := bufio.NewScanner( r )

  if skip, err := c.parseFile( scanner ); err != nil {
    return skip, err
  }

  if err := scanner.Err(); err != nil {
    return false, err
  }

  return true, nil
}

func (c *CIFImporter) parseFile( scanner *bufio.Scanner ) (bool, error) {

  var skip bool

  err := c.Update( func( tx *sql.Tx ) error {

    c.parserInit( tx )

    doImport, err := c.parseFileHeader( scanner )
    if err != nil {
      return err
    }

    if !doImport {
      skip = true
      return errors.New( "CIF too old" )
    }

    lastLine, err := c.parseTiplocs( scanner )
    if err != nil {
      return err
    }

    lastLine, err = c.parseAssociations( scanner, lastLine )
    if err != nil {
      return err
    }

    err = c.parseSchedules( scanner, lastLine )
    if err != nil {
      return err
    }

    return nil
  } );
  return skip, err
}

// Sets the CIF structure up to the current transaction
func (c *CIFImporter) parserInit( tx *sql.Tx ) {
  c.tx = tx
  c.curSchedule = nil
  c.update = false
}

// Looks for the initial HD record
// If the CIF is for a full import then reset the DB
func (c *CIFImporter) parseFileHeader( scanner *bufio.Scanner ) ( bool, error ) {
  if scanner.Scan() {

    line := scanner.Text()

    if line[0:2] == "HD" {

      doImport, err := c.parseHD( line )
      if err != nil {
        return false, err
      }

      return doImport, nil
    }
  }

  return false, errors.New( "Not a CIF file" )
}

// Parses the rest of the file after the header
func (c *CIFImporter) parseTiplocs( scanner *bufio.Scanner ) ( string, error ) {
  log.Println( "Parsing Tiploc's" )

  var line string

  for scanner.Scan() {
    line = scanner.Text()
    if bail, err := c.parseTiploc( line ); err != nil {
      return "", err
    } else if bail {
      return line, nil
    }
  }

  return line, nil
}

func (c *CIFImporter) parseTiploc( line string ) ( bool, error ) {
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
func (c *CIFImporter) parseAssociations( scanner *bufio.Scanner, lastLine string ) ( string, error ) {
  log.Println( "Parsing Associations" )

  line := lastLine

  // process the last line then continue & process that line
  if line != "" {
    if bail, err := c.parseAssociation( line ); err != nil {
      return "", err
    } else if bail {
      return line, nil
    }
  }

  for scanner.Scan() {
    line = scanner.Text()
    if bail, err := c.parseAssociation( line ); err != nil {
      return "", err
    } else if bail {
      return line, nil
    }
  }

  return line, nil
}

func (c *CIFImporter) parseAssociation( line string ) ( bool, error ) {
  switch line[0:2] {
    case "AA":
      return false, c.parseAA( line )

    // Ignore any unsupported records
    default:
      return true, nil
  }
}

// Parses the rest of the file after the header
func (c *CIFImporter) parseSchedules( scanner *bufio.Scanner, lastLine string ) error {
  log.Println( "Parsing Schedules" )

  count := 0

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
}

func (c *CIFImporter) parseSchedule( line string, count *int ) error {
  switch line[0:2] {
    case "BS":
      *count++
      if (*count % 5000) == 0 {
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
