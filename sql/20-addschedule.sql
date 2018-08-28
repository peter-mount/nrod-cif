-- ======================================================================
-- triggers to maintain the timetable
-- ======================================================================

-- ======================================================================
-- Insert/update a schedule
CREATE OR REPLACE FUNCTION timetable.addschedule( pSched JSON )
RETURNS BIGINT AS $$
DECLARE
  vsid    BIGINT;
  step    JSON;
  vord    SMALLINT;
  vstp    CHAR;
  sdt     DATE;
  edt     DATE;
  vdow    SMALLINT;
BEGIN
  vstp := pSched->'id'->>'stp';
  sdt := (pSched->'runs'->>'runsFrom')::DATE;
  edt := (pSched->'runs'->>'runsTo')::DATE;
  vdow := (pSched->'runs'->>'daysRun')::BIT(7)::INTEGER::SMALLINT;

  INSERT INTO timetable.schedule
    ( uid, stp, startdate, enddate, dow, entrydate )
    VALUES (
      pSched->'id'->>'uid',
      pSched->'id'->>'stp',
      sdt,
      edt,
      vdow,
      NOW()
    )
    ON CONFLICT ( uid, stp, startdate )
    DO UPDATE
      SET enddate = EXCLUDED.enddate,
          entrydate = EXCLUDED.entrydate
    RETURNING id INTO vsid;

  INSERT INTO timetable.schedule_json
    ( id, schedule )
    VALUES ( vsid, pSched )
    ON CONFLICT ( id )
    DO UPDATE
      SET schedule = pSched;

  -- Now remove & replace the station lookup
  DELETE FROM timetable.station WHERE sid = vsid;

  IF pSched->>'schedule' IS NOT NULL THEN
    vord := 0;
    FOR step IN SELECT * FROM json_array_elements( pSched->'schedule' )
    LOOP
      -- we only index against the public timetable
      IF step->'time'->>'pta' IS NOT NULL OR step->'time'->>'ptd' IS NOT NULL THEN
        INSERT INTO timetable.station
          ( sid, ord, tid, stp, startdate, enddate, dow, time )
          VALUES (
            vsid, vord, timetable.gettiplocid( step->> 'tpl' ), vstp,
            sdt, edt, vdow,
            (step->'time'->>'time')::TIME
          );
      END IF;
      vord := vord + 1;
    END LOOP;
  END IF;

  RETURN vsid;
END;
$$ LANGUAGE PLPGSQL;
