-- ==============================================================================================================
-- 20-schedule.sql contains utility functions used with the other schedule functions
-- ==============================================================================================================

-- ==============================================================================================================
-- Utility function that accepts either a schedule or an array of schedules in JSON format
-- and returns a map of the tiplocs within them, keyed by the tiploc code.
--
-- This map is usually returned to allow the client to display the correct name for a tiploc.
-- ==============================================================================================================

create or replace function timetable.getscheduletiplocs(pary json)
    returns json
as
$$
DECLARE
    tpls JSON;
BEGIN
    IF json_typeof(pary) = 'array' THEN
        -- Resolve tiplocs for each schedule in the array
        WITH sch AS (
            SELECT a -> 'schedule' schedule
            FROM json_array_elements(pary) a
            WHERE a ->> 'schedule' != 'null'
        ),
             locs AS (
                 SELECT json_array_elements(sched.schedule) loc
                 FROM (SELECT sch.schedule FROM sch) sched
             ),
             tiplocs AS (
                 SELECT DISTINCT loc ->> 'tpl' tpl
                 FROM locs
             )
        SELECT INTO tpls json_object_agg(tiploc, obj)
        FROM (
                 SELECT distinct on (t.tiploc) t.tiploc, row_to_json(t.*) as obj
                 FROM timetable.tiploc t
                 WHERE t.tiploc IN (
                     SELECT tpl
                     from tiplocs
                 )
                 ORDER BY t.tiploc
             ) tpl;
    ELSE
        -- Resolve tiplocs from a single schedule
        WITH tiplocs AS (
            SELECT DISTINCT loc ->> 'tpl' tpl
            FROM json_array_elements(pary -> 'schedule') loc
        )
        SELECT INTO tpls json_object_agg(tiploc, obj)
        FROM (
                 SELECT distinct on (t.tiploc) t.tiploc, row_to_json(t.*) as obj
                 FROM timetable.tiploc t
                 WHERE t.tiploc IN (
                     SELECT tpl
                     from tiplocs
                 )
                 ORDER BY t.tiploc
             ) tpl;
    END IF;


    RETURN tpls;
END ;
$$ LANGUAGE plpgsql;
