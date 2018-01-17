package cif

import (
  "fmt"
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "time"
)

// Tiploc represents a location on the rail network.
// This can be either a station, a junction or a specific point along the line/
type Tiploc struct {
  //XMLName         xml.Name  `xml:"tiploc"`
  // Tiploc key for this location
  Tiploc          string    `json:"tiploc" xml:"tiploc,attr"`
  NLC             int       `json:"nlc" xml:"nlc,attr"`
  NLCCheck        string    `json:"nlcCheck" xml:"nlcCheck,attr"`
  // Proper description for this location
  Desc            string    `json:"desc" xml:"desc,attr,omitempty"`
  // Stannox code, 0 means none
  Stanox          int       `json:"stanox" xml:"stanox,attr,omitempty"`
  // CRS code, "" for none. Codes starting with X or Z are usually not stations.
  CRS             string    `json:"crs" xml:"crs,attr,omitempty"`
  // NLC description of the location
  NLCDesc         string    `json:"nlcDesc" xml:"nlcDesc,attr,omitempty"`
  // The CIF extract this entry is from
  DateOfExtract   time.Time `json:"dateOfExtract" xml:"dateOfExtract,attr"`
  // Self (generated on rest only)
  Self            string    `json:"self,omitempty" xml:"self,attr,omitempty"`
}

// SetSelf sets the Self field to match this request
func (t *Tiploc) SetSelf( r *rest.Rest ) {
  t.Self = r.Self( "/tiploc/" + t.Tiploc )
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
func (c *CIF) TiplocHandler( r *rest.Rest ) error {
  return c.db.View( func( tx *bolt.Tx ) error {
    tpl := r.Var( "id" )

    response := &Response{}
    r.Value( response )

    if tiploc, exists := c.GetTiploc( tx, tpl ); exists {
      statistics.Incr( "tiploc.200" )
      r.Status( 200 )
      response.Status = 200
      response.Tiploc = []*Tiploc{tiploc}
      tiploc.SetSelf( r )
      response.Self = tiploc.Self
    } else {
      statistics.Incr( "tiploc.404" )
      r.Status( 404 )
      response.Status = 404
      response.Message = tpl + " not found"
    }

    return nil
  } )
}
