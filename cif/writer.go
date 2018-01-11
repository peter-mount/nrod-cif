package cif

import (
  "bytes"
  "encoding/binary"
//  "io"
//  "log"
  "time"
)

type BinaryCodec struct {
  buf  *bytes.Buffer
  err   error
}

func NewBinaryCodec() *BinaryCodec {
  var c *BinaryCodec = &BinaryCodec{}
  c.buf = &bytes.Buffer{}
  return c
}

func (c *BinaryCodec) Bytes() []byte {
  return c.buf.Bytes()
}

func (c *BinaryCodec) Error() error {
  return c.err
}

func (c *BinaryCodec) Write( i interface{ Write(*BinaryCodec) } ) *BinaryCodec {
  if c.err == nil {
    i.Write( c )
  }
  return c
}

func (c *BinaryCodec) WriteByte( b byte ) *BinaryCodec {
  if c.err == nil {
    c.err = c.buf.WriteByte( b )
  }
  return c
}

func (c *BinaryCodec) WriteBytes( b []byte ) *BinaryCodec {
  if c.err == nil {
    c.WriteInt32( int32( len( b ) ) )
  }
  if c.err == nil {
    _, c.err = c.buf.Write( b )
  }
  return c
}

func (c *BinaryCodec) WriteString( s string ) *BinaryCodec {
  if c.err == nil {
    c.WriteBytes( []byte( s ) )
  }
  return c
}

func (c *BinaryCodec) WriteInt( i int ) *BinaryCodec {
  return c.WriteInt64( int64(i) )
}

func (c *BinaryCodec) WriteInt64( i int64 ) *BinaryCodec {
  if c.err == nil {
    c.err = binary.Write( c.buf, binary.LittleEndian, i )
  }
  return c
}

func (c *BinaryCodec) WriteInt32( i int32 ) *BinaryCodec {
  if c.err == nil {
    c.err = binary.Write( c.buf, binary.LittleEndian, i )
  }
  return c
}

func (c *BinaryCodec) WriteTime( t time.Time ) *BinaryCodec {
  if c.err == nil {
    if b, err := t.MarshalBinary(); err != nil {
      c.err = err
    } else {
      _, c.err = c.buf.Write( b )
    }
  }
  return c
}
