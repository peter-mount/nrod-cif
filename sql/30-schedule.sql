-- ======================================================================
-- Return's a specific schedule.
--
-- ======================================================================

CREATE OR REPLACE FUNCTION timetable.timetableSchedule( pid TEXT )
RETURNS JSON AS $$
  WITH schedule AS (
    SELECT
        s.entrydate,
        j.schedule AS schedule
      FROM timetable.schedule s
        INNER JOIN timetable.schedule_json j on s.id = j.id
      WHERE s.id = id.decode( pid, 'radix.62' )
  ), locations AS (
    SELECT DISTINCT
        json_array_elements( (SELECT schedule->'schedule' FROM schedule) )->>'tpl' AS tpl
  ), tpls AS (
    -- View of all origin/destination tiplocs in the output
    SELECT DISTINCT
        t.*
      FROM timetable.tiploc t
        INNER JOIN locations l ON t.tiploc = l.tpl
      ORDER BY t.tiploc
  )
  SELECT json_build_object(
    -- The generated id for this schedule
    'id', pid,
    -- The schedule json
    'schedule', (SELECT s.schedule FROM schedule s LIMIT 1),
    -- Tiploc lookup map for entries within the schedule
    'tiploc', (SELECT json_object_agg(t.tiploc,row_to_json(t)) FROM tpls t ),
    -- The date this schedule was entered into the database
    'entrydate', (SELECT to_json(s.entrydate) FROM schedule s LIMIT 1),
    -- The timestamp of when this json was generated
    'generated', NOW() AT TIME ZONE 'Europe/London'
  )
$$ LANGUAGE SQL;