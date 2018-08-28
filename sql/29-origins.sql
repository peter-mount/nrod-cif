-- ======================================================================
-- Simple functions that return basic data about a schedule
--
-- timetable.origin( sid ) returns the origin tiploc and time
-- timetable.destination( sid ) returns the destination tiploc and time
--
-- ======================================================================
CREATE OR REPLACE FUNCTION timetable.origin( psid BIGINT )
RETURNS TABLE(
  sid     BIGINT,
  tiploc  VARCHAR,
  tod     TIME
) AS $$
  SELECT
      s.id,
      t.tiploc,
      s.time
    FROM timetable.tiploc t
    INNER JOIN timetable.station s ON t.id = s.tid
    WHERE s.sid = psid AND ord=0
    LIMIT 1
$$ LANGUAGE 'sql';

CREATE OR REPLACE FUNCTION timetable.destination( psid BIGINT )
RETURNS TABLE(
  sid     BIGINT,
  tiploc  VARCHAR,
  tod     TIME
) AS $$
  SELECT
      s.id,
      t.tiploc,
      s.time
    FROM timetable.tiploc t
    INNER JOIN timetable.station s ON t.id = s.tid
    WHERE s.sid = psid
    ORDER BY s.ord DESC
    LIMIT 1
$$ LANGUAGE 'sql';
