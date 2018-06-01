package cif

import (
  "fmt"
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/codec"
  "github.com/peter-mount/golib/rest"
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
  Desc            string    `json:"desc,omitempty" xml:"desc,attr,omitempty"`
  // Stannox code, 0 means none
  Stanox          int       `json:"stanox,omitempty" xml:"stanox,attr,omitempty"`
  // CRS code, "" for none. Codes starting with X or Z are usually not stations.
  CRS             string    `json:"crs,omitempty" xml:"crs,attr,omitempty"`
  // NLC description of the location
  NLCDesc         string    `json:"nlcDesc,omitempty" xml:"nlcDesc,attr,omitempty"`
  // True if this tiploc is a station
  Station         bool      `json:"station"`
  // The CIF extract this entry is from
  DateOfExtract   time.Time `json:"date" xml:"date,attr"`
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
  t.Update()
}

func (t *Tiploc) Update() {
  // Tiploc is a station IF it has a stanox, crs & crs not start with X or Z
  t.Station = t.Stanox > 0 &&t.CRS != "" && !(t.CRS[0] == 'X' || t.CRS[0] == 'Z')
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
