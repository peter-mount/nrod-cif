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
  sdt     DATE;
  edt     DATE;
BEGIN
  sdt := (pSched->'runs'->>'runsFrom')::DATE;
  edt := (pSched->'runs'->>'runsTo')::DATE;

  INSERT INTO timetable.schedule
    ( uid, stp, startdate, enddate, entrydate )
    VALUES (
      pSched->'id'->>'uid',
      pSched->'id'->>'stp',
      sdt,
      edt,
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
          ( sid, ord, tid, startdate, enddate, time )
          VALUES (
            vsid, vord, timetable.gettiplocid( step->> 'tpl' ),
            sdt, edt,
            (step->'time'->>'time')::TIME
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

DROP FUNCTION timetable.schedules( CHAR(3), TIMESTAMP with time zone, INTERVAL);

CREATE OR REPLACE FUNCTION timetable.schedules( pcrs CHAR(3), pst TIMESTAMP WITH TIME ZONE, dur INTERVAL )
RETURNS SETOF timetable.station AS $$
DECLARE
  ts TIMESTAMP WITHOUT TIME ZONE;
  sd DATE;
  ed DATE;
  st TIME;
  et TIME;
BEGIN
  -- Ensure we use the correct time of day during the summer
  ts = (pst AT TIME ZONE 'Europe/London'::TEXT)::TIMESTAMP WITHOUT TIME ZONE;
  sd = (ts::TEXT)::DATE;
  st = (ts::TEXT)::TIME;

  IF dur IS NULL OR dur < '0'::INTERVAL OR dur > '6 hours'::INTERVAL THEN
    ts = ts + '1 hour'::INTERVAL;
  ELSE
    ts = ts + dur;
  END IF;
  ed = (ts::TEXT)::DATE;
  et = (ts::TEXT)::TIME;

  IF st < et THEN
    RETURN QUERY
    WITH tpls AS (
      SELECT * FROM timetable.tiploc t
      WHERE stanox IN ( SELECT stanox FROM timetable.tiploc t2 WHERE crs = pCrs )
    )
    SELECT s.*
      FROM timetable.station s
      WHERE s.tid IN (SELECT id FROM tpls)
        AND s.time BETWEEN st AND et
        AND s.startdate <= sd
        AND s.enddate >=ed
      ORDER BY s.time, s.sid;
  ELSE
    RETURN QUERY
    WITH tpls AS (
      SELECT * FROM timetable.tiploc t
      WHERE stanox IN ( SELECT stanox FROM timetable.tiploc t2 WHERE crs = pCrs )
    )
    SELECT s.*
      FROM timetable.station s
      WHERE s.tid IN (SELECT id FROM tpls)
        AND (
          s.time BETWEEN st AND '24:00'::TIME
          OR s.time BETWEEN '00:00'::TIME AND et
        )
        AND s.startdate <= sd
        AND s.enddate >=sd
      ORDER BY s.time <= et, s.time, s.sid;
  END IF;
END;
$$ LANGUAGE PLPGSQL;

-- ======================================================================
