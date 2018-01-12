package cif

import (
  "encoding/json"
  "fmt"
  bolt "github.com/coreos/bbolt"
  "github.com/gorilla/mux"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/statistics"
  "log"
  "net/http"
  "time"
)

// Tiploc represents a location on the rail network.
// This can be either a station, a junction or a specific point along the line/
type Tiploc struct {
  // Tiploc key for this location
  Tiploc          string
  NLC             int
  NLCCheck        string
  // Proper description for this location
  Desc            string
  // Stannox code, 0 means none
  Stanox          int
  // CRS code, "" for none. Codes starting with X or Z are usually not stations.
  CRS             string
  // NLC description of the location
  NLCDesc         string
  // The CIF extract this entry is from
  DateOfExtract   time.Time
}

func (t *Tiploc) Write( c *codec.BinaryCodec ) {
  c.WriteString( t.Tiploc ).
    WriteInt( t.NLC ).
    WriteString( t.NLCCheck ).
    WriteString( t.Desc ).
    WriteInt( t.Stanox ).
    WriteString( t.CRS ).
    WriteString( t.NLCDesc ).
    WriteTime( t.DateOfExtract )
}

func (t *Tiploc) Read( c *codec.BinaryCodec ) {
  c.ReadString( &t.Tiploc ).
    ReadInt( &t.NLC ).
    ReadString( &t.NLCCheck ).
    ReadString( &t.Desc ).
    ReadInt( &t.Stanox ).
    ReadString( &t.CRS ).
    ReadString( &t.NLCDesc ).
    ReadTime( &t.DateOfExtract )
}

// String returns a human readable version of a Tiploc
func (t *Tiploc) String() string {
  return fmt.Sprintf(
    "Tiploc[%s, crs=%s, stanox=%05d, nlc=%d, desc=%s, nlcDesc=%s]",
    t.Tiploc,
    t.CRS,
    t.Stanox,
    t.NLC,
    t.Desc,
    t.NLCDesc )
}

// GetTiploc retrieves a Tiploc from the cif database
//
// tx An active readonly bolt.Tx
//
// t The Tiploc to retrieve, 1..7 characters long, always upper case
//
// Returns ( tiploc *Tiploc, exist bool )
//
// If exist is true then tiploc will be a new Tiploc instance with the retrieved data.
// If false then the tiploc is not in the database.
func (c *CIF) GetTiploc( tx *bolt.Tx, t string ) ( *Tiploc, bool ) {

  var tiploc *Tiploc = &Tiploc{}

  b := tx.Bucket( []byte("Tiploc") ).Get( []byte( t ) )
  if b == nil {
    return nil, false
  }

  codec.NewBinaryCodecFrom( b ).Read( tiploc )
  if tiploc.Tiploc == "" {
    return nil, false
  }

  return tiploc, true
}

// TiplocHandler implements a net/http handler that implements a simple Rest service to retrieve Tiploc records.
// The handler must have {id} set in the path for this to work, where id would represent the Tiploc code.
//
// For example:
//
// router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )
//
// where db is a pointer to an active CIF struct. When running this would allow GET requests like /tiploc/MSTONEE to return JSON representing that station.
func (c *CIF) TiplocHandler( w http.ResponseWriter, r *http.Request ) {
  var params = mux.Vars( r )

  tpl := params[ "id" ]

  if err := c.db.View( func( tx *bolt.Tx ) error {
    if tiploc, exists := c.GetTiploc( tx, tpl ); exists {
      statistics.Incr( "tiploc.200" )
      w.WriteHeader( 200 )
      json.NewEncoder( w ).Encode( tiploc )
    } else {
      statistics.Incr( "tiploc.404" )
      w.WriteHeader( 404 )
    }

    return nil
  } ); err != nil {
    statistics.Incr( "tiploc.500" )
    log.Println( "Get Tiploc", tpl, err )
    w.WriteHeader( 500 )
  }
}
