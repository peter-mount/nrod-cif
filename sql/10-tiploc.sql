-- ======================================================================
-- Functions & Triggers for managing tiplocs
-- ======================================================================

-- ======================================================================
-- Return's the unique ID for a tiploc
--
-- This algorithm encodes the characters 0-9, A-Z, a-z into a 6 bit number
-- and generates an integer constant based on the tiploc
--
-- ======================================================================
CREATE OR REPLACE FUNCTION timetable.gettiplocid(t varchar(7))
    RETURNS BIGINT AS
$$
DECLARE
    l CHAR;
    a INTEGER;
    c BIGINT;
BEGIN
    c = 0;
    FOR l IN SELECT regexp_split_to_table(t, '')
        LOOP
            a = ASCII(l);

            c = (c << 6) +
                CASE
                    WHEN a BETWEEN 48 AND 57 THEN a - 48
                    WHEN a BETWEEN 65 AND 90 THEN a - 55
                    WHEN a BETWEEN 97 AND 122 THEN a - 61
                    ELSE 0
                    END;
        END LOOP;
    RETURN c;
END;
$$ LANGUAGE PLPGSQL;

-- ======================================================================
-- Trigger to ensure that a tiploc's ID is correctly set
CREATE OR REPLACE FUNCTION timetable.settiplocid()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.id = timetable.gettiplocid(NEW.tiploc);
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;

CREATE TRIGGER settiploc
    BEFORE INSERT OR UPDATE
    ON timetable.tiploc
    FOR EACH ROW
EXECUTE PROCEDURE timetable.settiplocid();

-- ======================================================================
-- Fixes tiplocs so that those entries with no crs but the same stanox
-- get the valid crs code. This then allows us to ask for a timetable at
-- a crs and get all services.
--
-- For example VIC is VICTRIA but in the timetable they use either
-- VICTRIC or VICTRIE, VICTRIA has no services.
--
-- London Bridge is another example of this.
--
-- A valid crs code is one that does not start with X or Z
-- ======================================================================
create or replace function timetable.fixtiploccrs()
    returns void
as
$$
declare
    rec record;
begin
    for rec in select stanox, crs
               from timetable.tiploc
               where stanox > 0
                 and crs != ''
                 and substr(crs, 1, 1) not in ('X', 'Z')
        loop
            update timetable.tiploc
            set crs = rec.crs
            where crs = ''
              and stanox = rec.stanox;
        end loop;
end;
$$
    language plpgsql;