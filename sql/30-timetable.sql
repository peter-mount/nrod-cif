-- ======================================================================
-- Searches for schedules at a CRS for a specific timestamp
--
-- Unlike timetable.schedules this returns a more complex set needed for
-- the timetable tool at uktra.in/timetable
--
-- i.e. this includes origin, destination etc of this service
--
-- ======================================================================

CREATE OR REPLACE FUNCTION timetable.timetable(
  pcrs CHAR(3),
  pst TIMESTAMP WITH TIME ZONE
)
RETURNS JSON AS $$
	WITH schedules AS (
	  SELECT * FROM timetable.schedules( pcrs, date_trunc('hour',pst), null )
	), services AS (
	  SELECT
		  s.id AS sid,
		  s.uid, s.startDate, s.stp,
		  st.time AS "time",
		  ot.tiploc AS origin,
			ot.tod AS "originTime",
		  dt.tiploc AS destination,
			dt.tod AS "destinationTime"
		FROM timetable.schedule s
		  INNER JOIN schedules st ON st.sid = s.id
		  INNER JOIN timetable.origin( s.id ) ot ON s.id=ot.sid
		  INNER JOIN timetable.destination( s.id ) dt ON s.id = dt.sid
		ORDER BY st.time >= '01:00', st.time, s.id
  ), tpls AS (
    SELECT DISTINCT
        t.*
      FROM timetable.tiploc t
        INNER JOIN services s ON s.origin = t.tiploc OR s.destination = t.tiploc
      ORDER BY t.tiploc
  )
  SELECT json_build_object(
    -- tiploc entry for this station
    'station', (SELECT row_to_json(t) FROM timetable.tiploc t WHERE t.crs = pcrs LIMIT 1 ),
    -- timestamp for the start of the hour for this timetable
    'ts', (date_trunc('hour',pst) AT TIME ZONE 'Europe/London'::TEXT),
    -- Schedules within the timetable
    'schedules', (SELECT json_agg(row_to_json(s)) FROM services s ),
    -- tiploc entries for tiplocs within the timetable
    'tiploc', (SELECT json_agg(row_to_json(t)) FROM tpls t ),
    -- timestamp of when this timetable was generated
    'generated', (NOW() AT TIME ZONE 'Europe/London'::TEXT)
	)::JSON;
$$ LANGUAGE SQL;
