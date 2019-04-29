-- ======================================================================
-- Return's a specific schedule.
--
-- ======================================================================

CREATE OR REPLACE FUNCTION timetable.timetableSchedule(pid TEXT, pdate DATE)
    RETURNS JSON AS
$$
WITH schedule AS (
    SELECT s.id,
           s.uid,
           s.entrydate,
           j.schedule AS schedule
    FROM timetable.schedule s
             INNER JOIN timetable.schedule_json j ON s.id = j.id
    WHERE s.id = id.decode(pid, 'radix.62')
),
     assoc AS (
         -- Distinct on mainuid an cat
         SELECT DISTINCT ON (a.mainuid, a.cat) a.*,
                                               t.tiploc
         FROM timetable.assoc a
                  INNER JOIN schedule s
                             ON (a.mainuid = s.uid OR a.assocuid = s.uid)
                                 AND a.startdate <= pdate
                                 AND a.enddate >= pdate
                  INNER JOIN timetable.tiploc t ON t.id = a.tid
              -- order by end date then start date desc to get the one closest to t
         ORDER BY a.mainuid, a.cat, a.endDate, a.startDate DESC, a.stp
     ),
     assocSched AS (
         SELECT DISTINCT ON (s.uid) a.id AS aid,
                                    s.id AS sid,
                                    s.uid,
                                    s.startDate,
                                    s.stp
         FROM timetable.schedule s
                  INNER JOIN assoc a
                             ON (a.mainuid = s.uid OR a.assocuid = s.uid)
                                 AND s.startDate <= a.startDate
                                 AND s.endDate >= a.endDate
         WHERE s.id NOT IN ((SELECT id FROM schedule))
               -- order by uid & stp so that the distinct selects the correct one
               -- e.g. first stp in C, N, O P is selected handling overlays/cancellations correctly
         ORDER BY s.uid, s.startDate DESC, s.stp
     ),
     assocout AS (
         SELECT a.*,
                j.schedule
         FROM assoc a
                  INNER JOIN assocSched s ON a.id = s.aid
                  INNER JOIN timetable.schedule_json j ON j.id = s.sid
         ORDER BY a.cat
     ),
     locations AS (
         SELECT DISTINCT json_array_elements((SELECT schedule -> 'schedule' FROM schedule)) ->> 'tpl' AS tpl
         UNION
         SELECT DISTINCT
                 json_array_elements((SELECT schedule -> 'schedule' FROM assocout WHERE cat = 'JJ')) ->> 'tpl' AS tpl
         UNION
         SELECT DISTINCT
                 json_array_elements((SELECT schedule -> 'schedule' FROM assocout WHERE cat = 'NP')) ->> 'tpl' AS tpl
         UNION
         SELECT DISTINCT
                 json_array_elements((SELECT schedule -> 'schedule' FROM assocout WHERE cat = 'VV')) ->> 'tpl' AS tpl
     ),
     tpls AS (
         -- View of all origin/destination tiplocs in the output
         SELECT DISTINCT t.tiploc,
                         t.crs,
                         t.stanox,
                         t.name,
                         t.station
         FROM timetable.tiploc t
                  INNER JOIN locations l ON t.tiploc = l.tpl
         ORDER BY t.tiploc
     )
SELECT json_build_object(
           -- The generated id for this schedule
               'id', pid,
           -- The schedule json
               'schedule', (SELECT s.schedule FROM schedule s LIMIT 1),
               'associations', (SELECT json_agg(row_to_json(a)) FROM assocout a),
           -- Tiploc lookup map for entries within the schedule
               'tiploc', (SELECT json_object_agg(t.tiploc, row_to_json(t)) FROM tpls t),
           -- The date this schedule was entered into the database
               'entrydate', (SELECT to_json(s.entrydate) FROM schedule s LIMIT 1),
           -- The timestamp of when this json was generated
               'generated', NOW() AT TIME ZONE 'Europe/London'
           )
$$ LANGUAGE SQL;
