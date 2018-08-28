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
RETURNS TABLE (
  sid         BIGINT,
  uid         CHAR(6),
  startdate   DATE,
  stp         CHAR,
  tod         TIME,
  origin      VARCHAR,
  originTime  TIME,
  dest        VARCHAR,
  destTime    TIME
) AS $$
DECLARE
  ts    TIMESTAMP WITHOUT TIME ZONE;
  et    TIME;
BEGIN
  -- todo truncate pts to the hour
  pst = date_trunc('hour',pst);

  RETURN QUERY
    WITH schedules AS (
      SELECT * FROM timetable.schedules( pcrs, pst, null )
    )
    SELECT
        s.id::BIGINT,
        s.uid, s.startDate, s.stp,
        st.time,
        ot.tiploc, ot.tod,
        dt.tiploc, dt.tod
      FROM timetable.schedule s
        INNER JOIN schedules st ON st.sid = s.id
        INNER JOIN timetable.origin( s.id ) ot
        INNER JOIN timetable.destination( s.id ) dt
      ORDER BY st.time >= '01:00', st.time, s.id;
END;
$$ LANGUAGE PLPGSQL;

-- ======================================================================
