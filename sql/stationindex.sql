-- ======================================================================
-- triggers to maintain the timetable
-- ======================================================================

-- ======================================================================
-- Insert/update a schedule
CREATE OR REPLACE FUNCTION timetable.addschedule( pSched JSON )
RETURNS BIGINT AS $$
DECLARE
  vsid     BIGINT;
  step    JSON;
  vord     SMALLINT;
BEGIN
  INSERT INTO timetable.schedule
    ( uid, stp, startdate, enddate, entrydate )
    VALUES (
      pSched->'id'->>'uid',
      pSched->'id'->>'stp',
      (pSched->'runs'->>'runsFrom')::DATE,
      (pSched->'runs'->>'runsTo')::DATE,
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
          ( sid, ord, tid )
          VALUES (
            vsid,
            vord,
            timetable.gettiplocid( step->> 'tpl' )
          );
        vord := vord + 1;
      END IF;
    END LOOP;
  END IF;

  RETURN vsid;
END;
$$ LANGUAGE PLPGSQL;

-- ======================================================================
-- On schedule delete remove all entries
CREATE OR REPLACE FUNCTION timetable.scheddeleted()
RETURNS TRIGGER AS $$
BEGIN
  -- Delete the json
  DELETE FROM timetable.schedule_json sj
    WHERE sj.id = OLD.id;
  -- Delete the index
  DELETE FROM timetable.station s
    WHERE s.sid = OLD.id;

  RETURN OLD;
END;
$$ LANGUAGE PLPGSQL;

CREATE TRIGGER scheddeleted
  BEFORE DELETE ON timetable.schedule
  FOR EACH ROW
  EXECUTE PROCEDURE timetable.scheddeleted();

-- ======================================================================
