-- ======================================================================
-- Generates the unique scheduleId
--
-- This is formed of the number of days since the dateEpoch, the UID and the STP
-- formed into a single unique BIGINT
-- ======================================================================
CREATE OR REPLACE FUNCTION timetable.scheduleid(puid TEXT, psd DATE, pstp CHAR)
    RETURNS BIGINT AS
$$
DECLARE
    v BIGINT;
    a INTEGER;
    l CHAR;
BEGIN
    -- Start date first
    v = (FLOOR(EXTRACT(EPOCH FROM psd) / 86400) - id.meta_bigint('dateEpoch'))::BIGINT;

    -- Encode UID
    FOR l IN SELECT regexp_split_to_table(puid, '')
        LOOP
            a = ASCII(l);
            v = (v << 6) +
                CASE
                    WHEN a BETWEEN 48 AND 57 THEN a - 48
                    WHEN a BETWEEN 65 AND 90 THEN a - 55
                    WHEN a BETWEEN 97 AND 122 THEN a - 61
                    ELSE 0
                    END;
        END LOOP;

    -- Encode STP
    v = (v << 5) + (ASCII(pstp) - 65);

    RETURN v;
END;
$$ LANGUAGE PLPGSQL;

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
CREATE OR REPLACE FUNCTION timetable.dowmask(ts TIMESTAMP WITHOUT TIME ZONE)
    RETURNS INTEGER AS
$$
SELECT 1 << (7 - EXTRACT(ISODOW FROM ts)::INTEGER);
$$ LANGUAGE SQL;

-- ======================================================================
