package cifimport

import (
  "github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFImporter) parseAA( l string ) error {
  t := cif.Association{}

  // N, D or R
  action := l[2]

  i := 3
  i = parseString( l, i, 6, &t.MainUid )
  i = parseString( l, i, 6, &t.AssocUid )
  i = parseYYMMDD( l, i, &t.StartDate )
  i = parseYYMMDD( l, i, &t.EndDate )
  i = parseString( l, i, 7, &t.AssocDays )
  i = parseString( l, i, 2, &t.Category )
  i = parseString( l, i, 1, &t.DateInd )
  i = parseStringTrim( l, i, 7, &t.Tiploc )
  i = parseStringTrim( l, i, 1, &t.BaseSuffix )
  i = parseStringTrim( l, i, 1, &t.AssocSuffix )
  i += 1 // unused diagram type
  i = parseString( l, i, 1, &t.AssocType )
  i += 31
  i = parseString( l, i, 1, &t.STPIndicator )

  if action == 'D' {
    _, err := c.tx.Exec(
      "DELETE FROM timetable.assoc WHERE mainuid=$1 AND assocuid=$2 AND startdate=$3 AND stp=$4",
      t.MainUid,
      t.AssocUid,
      t.StartDate,
      t.STPIndicator,
    )
    return err
  }

  _, err := c.tx.Exec(
    "INSERT INTO timetable.assoc" +
    " (mainuid, assocuid, stp, startdate, enddate, dow, cat, dateInd, tid, baseSuffix, assocSuffix, assocType, entrydate )" +
    " VALUES (" +
    " $1, $2, $3, $4, $5," +
    // days of week
    " ($6)::BIT(7)::INTEGER::SMALLINT," +
    " $7, $8," +
    " timetable.gettiplocid( $9 )," +
    " $10, $11, $12," +
    " NOW()" +
    ") " +
    "ON CONFLICT ( mainuid, assocuid, stp, startdate ) " +
    "DO UPDATE SET " +
    "enddate = EXCLUDED.enddate," +
    "dow = EXCLUDED.dow," +
    "cat = EXCLUDED.cat," +
    "dateInd = EXCLUDED.dateInd," +
    "tid = EXCLUDED.tid," +
    "baseSuffix = EXCLUDED.baseSuffix," +
    "assocSuffix = EXCLUDED.assocSuffix," +
    "assocType = EXCLUDED.assocType," +
    "entrydate = EXCLUDED.entrydate",
    t.MainUid,
    t.AssocUid,
    t.STPIndicator,
    t.StartDate,
    t.EndDate,
    t.AssocDays,
    t.Category,
    t.DateInd,
    t.Tiploc,
    t.BaseSuffix,
    t.AssocSuffix,
    t.AssocType,
  )
  return err
}
