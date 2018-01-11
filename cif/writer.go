package cif

import (
  "bytes"
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
    c.WriteInt16( int16( len( b ) ) )
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

func (c *BinaryCodec) WriteStringArray( s []string ) *BinaryCodec {
  if c.err == nil {
    c.WriteInt16( int16( len( s ) ) )
    for _, v := range s {
      c.WriteString( v )
    }
  }
  return c
}

func (c *BinaryCodec) WriteInt( i int ) *BinaryCodec {
  return c.WriteInt64( int64(i) )
}

func (c *BinaryCodec) WriteInt64( v int64 ) *BinaryCodec {
  if c.err == nil {
    var b []byte = make( []byte, 8 )
    b[0] = byte(v)
  	b[1] = byte(v >> 8)
  	b[2] = byte(v >> 16)
  	b[3] = byte(v >> 24)
    b[4] = byte(v >> 32)
  	b[5] = byte(v >> 40)
  	b[6] = byte(v >> 48)
  	b[7] = byte(v >> 56)
    _, c.err = c.buf.Write( b )
  }
  return c
}

func (c *BinaryCodec) WriteInt32( v int32 ) *BinaryCodec {
  if c.err == nil {
    var b []byte = make( []byte, 4 )
    b[0] = byte(v)
  	b[1] = byte(v >> 8)
  	b[2] = byte(v >> 16)
  	b[3] = byte(v >> 24)
    _, c.err = c.buf.Write( b )
  }
  return c
}

func (c *BinaryCodec) WriteInt16( v int16 ) *BinaryCodec {
  if c.err == nil {
    var b []byte = make( []byte, 2 )
    b[0] = byte(v)
  	b[1] = byte(v >> 8)
    _, c.err = c.buf.Write( b )
  }
  return c
}


func (c *BinaryCodec) WriteTime( t time.Time ) *BinaryCodec {
  if c.err == nil {
    if b, err := t.MarshalBinary(); err != nil {
      c.err = err
    } else {
      c.WriteBytes( b )
    }
  }
  return c
}
