// CIF HD Record

package cif

import (
  "errors"
  "log"
  "time"
)

type HD struct {
  Id                      string    // 02 Record Identity, always "HD"
  FileMainframeIdentity   string    // 20
  DateOfExtract           time.Time // 06 Date DDMMYY 060315, 04 Time HHMM
  CurrentFileReference    string    // 07
  LastFileReference       string    // 07
  Update                  bool      // 01 U = Update = true, F = Full Extract = false
  Version                 string    // 01
  UserStartDate           time.Time // 06 DDMMYY
  UserEndDate             time.Time // 06 DDMMYY
  // Spare 20
}

// Parse HD record
// returns true if the file should be imported
func (c *CIF) parseHD( l string ) error {
  var h *HD = &HD{}

  i := 2
  i = parseString( l, i, 20, &h.FileMainframeIdentity )
  i = parseDDMMYY_HHMM( l, i, &h.DateOfExtract )
  i = parseString( l, i, 7, &h.CurrentFileReference )
  i = parseString( l, i, 7, &h.LastFileReference )

  var update string
  i = parseString( l, i, 1, &update )
  h.Update = update == "U"

  i = parseString( l, i, 1, &h.Version )
  i = parseDDMMYY( l, i, &h.UserStartDate )
  i = parseDDMMYY( l, i, &h.UserEndDate )

  log.Printf(
    "CIF %s Extracted %v Date Range %v - %v\n",
    h.FileMainframeIdentity,
    h.DateOfExtract.Format( HumanDateTime ),
    h.UserStartDate.Format( HumanDate ),
    h.UserEndDate.Format( HumanDate ) )

  if h.Update {
    // Check existing to see we are more recent, skip file if before
    if c.Header != nil && h.UserStartDate.After( c.Header.UserStartDate ) {
      log.Println( "File is too old" )
      return errors.New( "CIF File is too old" )
    }

    log.Println( "Performing CIF Update" )
  } else {
    log.Println( "Performing Full import" )
  }

  c.Header = h

  return nil
}
