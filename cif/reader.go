package cif

import (
  "bytes"
  "encoding/binary"
  "time"
)

func NewBinaryCodecFrom( b []byte ) *BinaryCodec {
  var c *BinaryCodec = &BinaryCodec{}
  buf := bytes.NewBuffer( b )
  c.buf = buf
  return c
}

func (c *BinaryCodec) Read( i interface{ Read(*BinaryCodec) } ) *BinaryCodec {
  if c.err == nil {
    i.Read( c )
  }
  return c
}

func (c *BinaryCodec) ReadByte( b *byte ) *BinaryCodec {
  if c.err == nil {
    if t, err := c.buf.ReadByte(); err != nil {
      c.err = err
    } else {
      *b = t
    }
  }
  return c
}

func (c *BinaryCodec) ReadBytes( b *[]byte ) *BinaryCodec {
  var l int32
  if c.err == nil {
    c.ReadInt32( &l )
  }
  if c.err == nil {
    ary := make( []byte, l )
    if _, err := c.buf.Read( ary ); err != nil {
      c.err = err
    } else {
      *b = ary
    }
  }
  return c
}

func (c *BinaryCodec) ReadString( s *string) *BinaryCodec {
  if c.err == nil {
    var b []byte
    c.ReadBytes( &b )
    if c.err == nil {
      *s = string(b[:])
    }
  }
  return c
}

func (c *BinaryCodec) ReadInt( i *int ) *BinaryCodec {
  var t int64
  c.ReadInt64( &t )
  if( c.err == nil ) {
    *i = int(t)
  }
  return c
}

func (c *BinaryCodec) ReadInt64( i *int64) *BinaryCodec {
  if c.err == nil {
    c.err = binary.Read( c.buf, binary.LittleEndian, i )
  }
  return c
}

func (c *BinaryCodec) ReadInt32( i *int32 ) *BinaryCodec {
  var t int32
  if c.err == nil {
    c.err = binary.Read( c.buf, binary.LittleEndian, &t )
  }
  if c.err == nil {
    *i = t
  }
  return c
}

func (c *BinaryCodec) ReadTime( i *time.Time ) *BinaryCodec {
  var b []byte

  if c.err == nil {
    c.ReadBytes( &b )
  }

  if c.err == nil {
    c.err = i.UnmarshalBinary( b )
  }

  return c
}
