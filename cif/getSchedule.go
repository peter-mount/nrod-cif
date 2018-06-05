package cif

import (
  "database/sql"
  "encoding/json"
  "time"
)

// GetSchedule returns a specific Schedule's for a specific TrainUID, startDate and STPIndicator
// If no schedule exists for the required key then nil is returned
func (c *CIF) GetSchedule( uid string, date time.Time, stp string ) (*Schedule, error) {
  row := c.db.QueryRow(
    "SELECT j.schedule FROM timetable.schedule_json j INNER JOIN timetable.schedule s ON j.id=s.id WHERE s.uid=$1 AND s.startdate=$2 AND s.stp=$3",
    uid, date, stp )
  j := ""
  err := row.Scan( &j )
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }

  s := &Schedule{}
  err = json.Unmarshal( []byte(j), s )
  if err != nil {
    return nil, err
  }

  return s, nil
}

func (c *CIF) GetScheduleQuery( join string, where string, args ...interface{} ) ([]*Schedule, error) {
  rows, err := c.db.Query(
    "SELECT j.schedule FROM timetable.schedule_json j INNER JOIN timetable.schedule s ON j.id=s.id " + join + " WHERE " + where + " ORDER BY uid, startdate, stp",
    args... )
  if err == sql.ErrNoRows {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var res []*Schedule
  for rows.Next() {
    var j string
    err := rows.Scan( &j )
    if err != nil {
      return nil, err
    }

    if j != "" {
      s := &Schedule{}
      err = json.Unmarshal( []byte(j), s )
      if err != nil {
        return nil, err
      }
      res = append( res, s )
    }
  }

  return res, nil
}

// GetSchedulesByUID returns a slice of all available schedules for a UID
func (c *CIF) GetSchedulesByUID( uid string ) ([]*Schedule, error) {
  s, err := c.GetScheduleQuery( "", "uid=$1", uid )
  return s, err
}
