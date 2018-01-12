// CIF HD Record

package cif

import (
  bolt "github.com/coreos/bbolt"
  "errors"
  "fmt"
  "github.com/peter-mount/golib/codec"
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

func (h *HD) Write( c *codec.BinaryCodec ) {
  c.WriteString( h.Id ).
    WriteString( h.FileMainframeIdentity ).
    WriteTime( h.DateOfExtract ).
    WriteString( h.CurrentFileReference ).
    WriteString( h.LastFileReference ).
    WriteBool( h.Update ).
    WriteString( h.Version ).
    WriteTime( h.UserStartDate ).
    WriteTime( h.UserEndDate )
}

func (h *HD) Read( c *codec.BinaryCodec ) {
  c.ReadString( &h.Id ).
    ReadString( &h.FileMainframeIdentity ).
    ReadTime( &h.DateOfExtract ).
    ReadString( &h.CurrentFileReference ).
    ReadString( &h.LastFileReference ).
    ReadBool( &h.Update ).
    ReadString( &h.Version ).
    ReadTime( &h.UserStartDate ).
    ReadTime( &h.UserEndDate )
}

func (c *CIF) GetHD() ( *HD, error ) {
  var h *HD = &HD{}

  if err := c.db.View( func( tx *bolt.Tx) error {

    b := tx.Bucket( []byte("Meta") ).Get( []byte( "lastCif" ) )
    if b != nil {
      codec.NewBinaryCodecFrom( b ).Read( h )
    }
    return nil
  }); err != nil {
    return nil, err
  }

  return h, nil
}

// Parse HD record
// returns true if the file should be imported
func (c *CIF) parseHD( l string ) error {
  var h *HD = &HD{}

  i := 0
  i = parseString( l, i, 2, &h.Id )
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

  log.Println( h.String() )

  if h.Update {
    // Check existing to see we are more recent, skip file if before
    if c.header != nil && ( c.header.UserStartDate.After( h.UserStartDate ) || c.header.UserStartDate.Equal( h.UserStartDate ) ) {
      log.Println( "File is too old" )
      return errors.New( "CIF File is too old" )
    }

    log.Println( "Performing CIF Update" )
  } else {
    log.Println( "Performing Full import" )
  }

  c.importhd = h
  return nil
}

func (h *HD ) String() string {
  return fmt.Sprintf(
    "CIF %s Extracted %v Date Range %v - %v Update %v",
    h.FileMainframeIdentity,
    h.DateOfExtract.Format( HumanDateTime ),
    h.UserStartDate.Format( HumanDate ),
    h.UserEndDate.Format( HumanDate ),
    h.Update )
}
