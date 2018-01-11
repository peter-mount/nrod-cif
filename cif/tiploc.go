package cif

import (
  "encoding/json"
  "fmt"
  bolt "github.com/coreos/bbolt"
  "github.com/gorilla/mux"
  "github.com/peter-mount/golib/statistics"
  "log"
  "net/http"
  "time"
)

type Tiploc struct {
  Tiploc          string
  NLC             int
  NLCCheck        string
  Desc            string
  Stanox          int
  CRS             string
  NLCDesc         string
  // The CIF extract this entry is from
  DateOfExtract   time.Time
}

func (t *Tiploc) Write( c *BinaryCodec ) {
  c.WriteString( t.Tiploc ).
    WriteInt( t.NLC ).
    WriteString( t.NLCCheck ).
    WriteString( t.Desc ).
    WriteInt( t.Stanox ).
    WriteString( t.CRS ).
    WriteString( t.NLCDesc ).
    WriteTime( t.DateOfExtract )
}

func (t *Tiploc) Read( c *BinaryCodec ) {
  c.ReadString( &t.Tiploc ).
    ReadInt( &t.NLC ).
    ReadString( &t.NLCCheck ).
    ReadString( &t.Desc ).
    ReadInt( &t.Stanox ).
    ReadString( &t.CRS ).
    ReadString( &t.NLCDesc ).
    ReadTime( &t.DateOfExtract )
}

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

func (c *CIF) GetTiploc( tx *bolt.Tx, t string ) ( *Tiploc, bool ) {

  var tiploc *Tiploc = &Tiploc{}

  b := tx.Bucket( []byte("Tiploc") ).Get( []byte( t ) )
  if b == nil {
    return nil, false
  }

  NewBinaryCodecFrom( b ).Read( tiploc )
  if tiploc.Tiploc == "" {
    return nil, false
  }

  return tiploc, true
}

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
