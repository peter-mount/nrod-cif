package cif

import (
  bolt "github.com/coreos/bbolt"
  "github.com/peter-mount/golib/rest"
  "fmt"
  "strconv"
)

func (c *CIF) TimetableHandler( r *rest.Rest ) error {
  return c.db.View( func( tx *bolt.Tx ) error {

    crs := r.Var( "crs" )
    date := r.Var( "date" )
    hour, _ := strconv.Atoi( r.Var( "hour" ) )

    //bucket := tx.Bucket( []byte( "Schedule" ) )
    //cursor := bucket.Cursor()

    result := NewResponse()

    result.SetSelf( r, fmt.Sprintf( "/timetable/%s/%s/%d", crs, date, hour ) )

    return nil
  })
}
