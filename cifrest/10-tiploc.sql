-- ==============================================================================================================
-- 10-tiploc.sql contains rest endpoints that return tiploc's or set's of tiplocs.
--
-- Specifically:
-- /ref/tiploc/{tiploc}     returns the named tiploc
-- /ref/crs/{crs}           returns an array of tiplocs for a specific crs
-- /ref/stanox/{stanox}     returns an array of tiplocs for a specific stanox
-- ==============================================================================================================

-- ==============================================================================================================
-- getcrs - handles the /cif/{crs} GET endpoint

create or replace function timetable.getcrs(pcrs text)
    returns json
as
$$
DECLARE
    tpls JSON;
    name TEXT;
BEGIN

    SELECT INTO tpls json_agg(tpl)
    FROM (
             SELECT *
             FROM timetable.tiploc
             WHERE crs = pcrs
             ORDER BY nlcdesc DESC, stanox, nlc
         ) tpl;

    SELECT INTO name a ->> 'name' from (select json_array_elements(tpls) as a limit 1) a;

    RETURN json_build_object(
            'crs', pcrs,
            'name', name,
            'tiploc', tpls
        );

END;
$$
    language plpgsql;

-- ==============================================================================================================
-- getstanox returns all tiplocs with a specific stanox
create or replace function timetable.getstanox(pstanox int)
    returns json
as
$$
DECLARE
    tpls JSON;
BEGIN

    SELECT INTO tpls json_agg(tpl)
    FROM (
             SELECT *
             FROM timetable.tiploc
             WHERE stanox = pstanox
             ORDER BY id
         ) tpl;

    RETURN tpls;

END;
$$
    language plpgsql;

-- ==============================================================================================================
-- gettiploc returns details of a single tiploc
create or replace function timetable.gettiploc(ptpl text)
    returns json
as
$$
DECLARE
    tpls JSON;
BEGIN

    SELECT INTO tpls row_to_json(t.*)
    FROM timetable.tiploc t
    WHERE tiploc = ptpl;

    RETURN tpls;

END;
$$
    language plpgsql;
