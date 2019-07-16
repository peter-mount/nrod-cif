-- ==============================================================================================================
-- 21-schedule.sql contains functions to return a raw schedule
--
-- /schedule/{uid}/{date}/{stp} returns a specific schedule
-- /schedule/{uid}              returns an array of schedules for a specific uid
-- ==============================================================================================================

-- ==============================================================================================================
-- Returns the details of a specific schedule

create or replace function timetable.getschedule(puid text, pdate date, pstp text)
    returns json
as
$$
DECLARE
    sched JSON;
    tpls  JSON;
BEGIN

    -- Locate the schedule
    SELECT INTO sched j.schedule
    FROM timetable.schedule_json j
             INNER JOIN timetable.schedule s ON j.id = s.id
    WHERE s.uid = puid
      AND s.startdate = pdate
      AND s.stp = pstp;

    -- Bail if not found
    IF NOT FOUND THEN
        RETURN NULL;
    end if;

    -- Resolve the tiplocs from the schedule
    tpls = timetable.getscheduletiplocs(sched);

    RETURN json_build_object(
            'uid', puid,
            'date', pdate,
            'stp', pstp,
            'schedule', sched,
            'tiploc', tpls
        );

END;
$$
    language plpgsql;

-- ==============================================================================================================
-- Returns the details of all schedules with a specific uid

create or replace function timetable.getschedule(puid text)
    returns json
as
$$
DECLARE
    sched JSON;
    tpls  JSON;
BEGIN

    SELECT INTO sched json_agg(s1)
    FROM (
             SELECT j.schedule as s1
             FROM timetable.schedule_json j
                      INNER JOIN timetable.schedule s ON j.id = s.id
             WHERE s.uid = puid
         ) s0;

    IF NOT FOUND THEN
        RETURN NULL;
    end if;

    tpls = timetable.getscheduletiplocs(sched);

    RETURN json_build_object(
            'uid', puid,
            'schedules', sched,
            'tiploc', tpls
        );

END ;
$$
    language plpgsql;

-- ==============================================================================================================
