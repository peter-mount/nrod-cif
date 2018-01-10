package cif

import (
  //"log"
  "strconv"
  "strings"
  "time"
)

const (
  DateTime        = "2006-01-02 15:04:05"
  Date            = "2006-01-02"
  HumanDateTime   = "2006 Jan 02 15:04:05"
  HumanDate       = "2006 Jan 02"
  Time            = "15:04:05"
)

func parseString( line string, s int, l int, v *string ) int {
  *v = line[s:s+l]
  return s + l
}

func parseStringTrim( line string, s int, l int, v *string ) int {
  var st string
  var ret = parseString( line, s, l, &st )
  *v = strings.Trim( st, " " )
  return ret
}

// Parse a string, trim then title it
func parseStringTitle( line string, s int, l int, v *string ) int {
  var st string
  var ret = parseStringTrim( line, s, l, &st )
  *v = strings.Title( strings.ToLower( st ) )
  return ret
}

// Parse DDMMYY to Time
func parseDDMMYY( line string, s int, v *time.Time ) int {
  var dt string
  var ret = parseString( line, s, 6, &dt )
  t, _ := time.Parse( "20060102", "20" + dt[4:6] + dt[2:4] + dt[0:2] )
  *v = t
  return ret
}

// Parse YYMMDD to Time
func parseYYMMDD( line string, s int, v *time.Time ) int {
  var dt string
  var ret = parseString( line, s, 6, &dt )
  t, _ := time.Parse( "20060102", "20" + dt )
  *v = t
  return ret
}

// Parse DDMMYY & HHMM to Time
func parseDDMMYY_HHMM( line string, s int, v *time.Time ) int {
  var dt string
  var ret = parseString( line, s, 10, &dt )
  t, _ := time.Parse( "200601021504", "20" + dt[4:6] + dt[2:4] + dt[0:2] + dt[6:] )
  *v = t
  return ret
}

func parseInt( line string, s int, l int, v *int ) int {
  var i string
  var ret = parseString( line, s, l, &i )
  val, _ := strconv.Atoi( i )
  *v = val
  return ret
}
