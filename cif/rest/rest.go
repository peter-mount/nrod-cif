package cif

import (
  "github.com/peter-mount/golib/rest"
)

func (c *CIF) InitRest( ctx *rest.ServerContext ) {
  ctx.Handle( "/crs/{id}", c.CRSHandler ).Methods( "GET" )
  ctx.Handle( "/stanox/{id}", c.StanoxHandler ).Methods( "GET" )
  ctx.Handle( "/tiploc/{id}", c.TiplocHandler ).Methods( "GET" )

  ctx.Handle( "/schedule/{uid}/{date}/{stp}", c.ScheduleHandler ).Methods( "GET" )
  ctx.Handle( "/schedule/{uid}", c.ScheduleUIDHandler ).Methods( "GET" )

  ctx.Handle( "/timetable/{crs}/{date}/{hour}", c.TimetableHandler ).Methods( "GET" )

  ctx.Handle( "/importCIF", c.ImportHandler ).Methods( "POST" )
}
