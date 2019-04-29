package cifimport

import (
  "github.com/peter-mount/nre-feeds/util"
  "strconv"
  "strings"
  "time"
)

const (
  DateTime      = "2006-01-02 15:04:05"
  Date          = "2006-01-02"
  HumanDateTime = "2006 Jan 02 15:04:05"
  HumanDate     = "2006 Jan 02"
  Time          = "15:04:05"
)

func parseString(line string, s int, l int, v *string) int {
  *v = line[s : s+l]
  return s + l
}

func parseStringTrim(line string, s int, l int, v *string) int {
  var st string
  var ret = parseString(line, s, l, &st)
  *v = strings.Trim(st, " ")
  return ret
}

// Parse a string, trim then title it
func parseStringTitle(line string, s int, l int, v *string) int {
  var st string
  var ret = parseStringTrim(line, s, l, &st)
  *v = strings.Title(strings.ToLower(st))
  return ret
}

// Parse DDMMYY to Time
func parseDDMMYY(line string, s int, v *time.Time) int {
  var dt string
  var ret = parseString(line, s, 6, &dt)
  t, _ := time.Parse("20060102", "20"+dt[4:6]+dt[2:4]+dt[0:2])
  *v = t
  return ret
}

// Parse YYMMDD to Time
func parseYYMMDD(line string, s int, v *time.Time) int {
  var dt string
  var ret = parseString(line, s, 6, &dt)
  t, _ := time.Parse("20060102", "20"+dt)
  *v = t
  return ret
}

// Parse DDMMYY & HHMM to Time
func parseDDMMYY_HHMM(line string, s int, v *time.Time) int {
  var dt string
  var ret = parseString(line, s, 10, &dt)
  t, _ := time.Parse("200601021504", "20"+dt[4:6]+dt[2:4]+dt[0:2]+dt[6:])
  *v = t
  return ret
}

func parseInt(line string, s int, l int, v *int) int {
  var i string
  var ret = parseString(line, s, l, &i)
  val, _ := strconv.Atoi(i)
  *v = val
  return ret
}

// Parse HHMM into time of day in seconds, -1 if none
func parseHHMM(l string, s int, v **util.PublicTime) int {
  var a, b int
  var c string
  var ret = parseString(l, s, 4, &c)
  if c == "    " {
    *v = nil
  } else {
    *v = &util.PublicTime{}
    a, _ = strconv.Atoi(c[0:2])
    b, _ = strconv.Atoi(c[2:4])
    (**v).Set((a * 60) + b)
  }
  return ret
}

// Parse HHMMS into time of day in seconds, -1 if none. S is "H" for 30 seconds past minute
func parseHHMMS(l string, s int, v **util.WorkingTime) int {
  var a *util.PublicTime
  var b string
  var ret = parseHHMM(l, s, &a)
  ret = parseString(l, ret, 1, &b)
  if a == nil || a.IsZero() {
    *v = nil
  } else {
    *v = &util.WorkingTime{}
    if b == "H" {
      (**v).Set((a.Get() * 60) + 30)
    } else {
      (**v).Set(a.Get() * 60)
    }
  }
  return ret
}

func parseActivity(l string, s int, v *[]string) int {
  var ary []string
  i := s
  for j := 0; j < 6; j++ {
    var act string
    i = parseStringTrim(l, i, 2, &act)
    if act != "" {
      ary = append(ary, act)
    }
  }
  *v = ary
  return i
}
