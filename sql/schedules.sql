-- ======================================================================
-- Return a bitmask for the day of week a service is running on
--
-- In CIF the dow is a 7 bit mask, Monday being first entry & Sunday last
--
-- ISODOW is Mon=1, Sun=7 so bitmask is 1<<(7-ISODOW)
--
-- so that for Monday we have 1<<6 and for sunday 1<<0
--
-- Common values:
--
--  Weekdays    124
--  Weekend       3
--  Every day   127
--
--  Monday       64
--  Tuesday      32
--  Wednesday    16
--  Thursday      8
--  Friday        4
--  Saturday      2
--  Sunday        1
--
-- ======================================================================
CREATE OR REPLACE FUNCTION timetable.dowmask( ts TIMESTAMP WITHOUT TIME ZONE )
RETURNS INTEGER AS $$
  SELECT 1 << (7 - EXTRACT( ISODOW FROM ts )::INTEGER);
$$ LANGUAGE SQL;

-- ======================================================================
-- Searches for schedules at a CRS for a specific timestamp
--
-- pcrs The CRS code of the station.
-- pst  The timestamp of the start of the range required.
-- dur  Interval describing how long to search,
--      null for the default of 1 hour.
--
-- If dur is negative or > 6 hours then an error will occur
--
-- ======================================================================

CREATE OR REPLACE FUNCTION timetable.schedules(
  pcrs CHAR(3),
  pst TIMESTAMP WITH TIME ZONE,
  dur INTERVAL
)
RETURNS SETOF timetable.station AS $$
DECLARE
  ts    TIMESTAMP WITHOUT TIME ZONE;
  sd    DATE;
  ed    DATE;
  st    TIME;
  et    TIME;
  sdow  INTEGER;
  edow  INTEGER;
BEGIN
  -- Ensure we use the correct time of day during the summer
  ts = (pst AT TIME ZONE 'Europe/London'::TEXT)::TIMESTAMP WITHOUT TIME ZONE;
  sd = (ts::TEXT)::DATE;
  st = (ts::TEXT)::TIME;
  sdow = timetable.dowmask( ts );

  IF dur IS NULL THEN
    ts = ts + '1 hour'::INTERVAL;
  ELSIF dur < '1 minute'::INTERVAL OR dur > '6 hours'::INTERVAL THEN
    RAISE EXCEPTION 'Invalid search range %s', dur
      USING HINT = 'Keep range between 1 minute and 6 hours';
  ELSE
    ts = ts + dur;
  END IF;
  ed = (ts::TEXT)::DATE;
  et = (ts::TEXT)::TIME;
  edow = timetable.dowmask( ts );

  IF st < et THEN
    -- Range is in the same day - the usual senario
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
        -- just test sdow as edow will be the same in this instance
        AND (s.dow & sdow) = sdow
      ORDER BY s.time, s.sid;
  ELSE
    -- Range crosses midnight so 2 searches, one before and one after midnight
    -- as obviously the dates and dow's are different
    RETURN QUERY
    WITH tpls AS (
      SELECT * FROM timetable.tiploc t
      WHERE stanox IN ( SELECT stanox FROM timetable.tiploc t2 WHERE crs = pCrs )
    ),
    merged AS (
      SELECT s.*
        FROM timetable.station s
        WHERE s.tid IN (SELECT id FROM tpls)
          AND s.time BETWEEN st AND '24:00'::TIME
          AND s.startdate <= sd
          AND s.enddate >=sd
          AND (s.dow & sdow) = sdow
      UNION
      SELECT s.*
        FROM timetable.station s
        WHERE s.tid IN (SELECT id FROM tpls)
          AND s.time BETWEEN '00:00'::TIME AND et
          AND s.startdate <= ed
          AND s.enddate >=ed
          AND (s.dow & edow) = edow
    )
    SELECT s.* from merged s
      -- order by time & then sid.
      -- s.time <= et ensures we sort the times correctly with those past
      -- midnight after those before midnight
      ORDER BY s.time <= et, s.time, s.sid;
  END IF;
END;
$$ LANGUAGE PLPGSQL;

-- ======================================================================
