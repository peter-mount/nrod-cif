-- ================================================================================
-- Various ID generators
--
-- Instagram: returns a unique id thats shardable
--
-- Sources:
--
-- Instagram        http://rob.conery.io/2014/05/29/a-better-id-generator-for-postgresql/
-- pseudo encrypt   http://wiki.postgresql.org/wiki/Pseudo_encrypt
--
-- ================================================================================

CREATE SCHEMA IF NOT EXISTS id;

-- ================================================================================
-- Misc global data for generators
CREATE TABLE IF NOT EXISTS id.meta
(
    -- Name of parameter
    name  NAME NOT NULL,
    -- Value(s), one each of different types
    ival  INTEGER,
    bval  BIGINT,
    tval  TEXT,
    -- Optional description of what param is for
    descr TEXT,
    PRIMARY KEY (name)
);

-- Return the named integer value or null
CREATE OR REPLACE FUNCTION id.meta_int(_name NAME)
    RETURNS INTEGER
AS
$$
SELECT ival
FROM id.meta
WHERE name = _name;
$$ LANGUAGE SQL
   IMMUTABLE;

-- Return the named bigint value or null
CREATE OR REPLACE FUNCTION id.meta_bigint(_name NAME)
    RETURNS BIGINT
AS
$$
SELECT bval
FROM id.meta
WHERE name = _name;
$$ LANGUAGE SQL;

-- Return the named text value or null
CREATE OR REPLACE FUNCTION id.meta_text(_name NAME)
    RETURNS TEXT
AS
$$
SELECT tval
FROM id.meta
WHERE name = _name;
$$ LANGUAGE SQL;

-- ================================================================================
--
-- Global config. This will setup defaults for shardId and epoch. If installing
-- in a distributed system or when used within a shared system ensure you override
-- these to the common values for your setup.
--
-- shardId  The database shard, either -1 for no sharding or 0..7 for a specific instance
-- epoch    The epoch in milliseconds to base ID generation on.
--
-- ================================================================================
--
-- The DB shard, must be set for each schema shard you have.
-- Valid values are 0.. 8191 (1<10)
INSERT INTO id.meta (name, ival, descr)
VALUES ('shardId', 1, 'Database shard id, unique per instance')
ON CONFLICT DO NOTHING;

-- The epoch, usually the point in time to start from - must be common to all shards
-- Original 1314220021721 but 1489586681104 makes smaller id's and valid for 2017-04-18
INSERT INTO id.meta (name, bval, descr)
VALUES ('epoch', 1492501747000, 'Epoch (milliseconds) for base of generated ID''s')
ON CONFLICT DO NOTHING;

-- The date epoch, used by some date calcs for smaller id's, = 2017-04-18
INSERT INTO id.meta (name, bval, descr)
VALUES ('dateEpoch', FLOOR(EXTRACT(EPOCH FROM '2017-04-18'::date) / 86400), 'Epoch (days) for base of generated ID''s')
ON CONFLICT DO NOTHING;

-- Alternate to when this file was inserted.
--INSERT INTO id.meta (name,bval)
--    VALUES ( 'epoch', FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) )
--    ON CONFLICT DO NOTHING;


-- ================================================================================
-- The instagram id generator
--
-- See https://engineering.instagram.com/sharding-ids-at-instagram-1cf5a71e5a5c

CREATE SEQUENCE IF NOT EXISTS id.instagram;

-- Create an id for a specific shard
CREATE OR REPLACE FUNCTION id.instagram(shard_id INTEGER, OUT result BIGINT)
AS
$$
DECLARE
    our_epoch  BIGINT := id.meta_bigint('epoch');
    seq_id     BIGINT;
    now_millis BIGINT;
BEGIN
    SELECT nextval('id.instagram') % 1024 INTO seq_id;

    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;

    result := (now_millis - our_epoch) << 23;
    result := result | (shard_id << 10);
    result := result | (seq_id);
END;
$$ LANGUAGE PLPGSQL;

-- Create an ID for this shard
CREATE OR REPLACE FUNCTION id.instagram()
    RETURNS BIGINT
AS
$$
SELECT id.instagram(id.meta_int('shardId'))
$$ LANGUAGE SQL;

-- ================================================================================
-- pseudo_encrypt a simple one way function that can simply hide the fact that
-- data has been generated by an INTEGER sequence.

CREATE OR REPLACE FUNCTION id.pseudo_encrypt(VALUE INTEGER)
    RETURNS INTEGER
AS
$$
DECLARE
    l1 int;
    l2 int;
    r1 int;
    r2 int;
    i  int := 0;
BEGIN
    l1 := (VALUE >> 16) & 65535;
    r1 := VALUE & 65535;
    WHILE i < 3
        LOOP
            l2 := r1;
            r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * 32767)::int;
            l1 := l2;
            r1 := r2;
            i := i + 1;
        END LOOP;
    RETURN ((r1 << 16) + l1);
END;
$$ LANGUAGE plpgsql
   strict
   immutable;

-- ================================================================================
-- Encode a bigint into a string using a specific radix
--
-- Useage:
--      id.encode( val )                Generate a plain radix.36 value
--      id.encode( val, 'radix.36' )    Generate using the named alphabet
--
-- Predefined alphabets:
--      hex.u           Radix 16 of 0-9 A-F         Also known as hex
--      hex.l           Radix 16 of 0-9 a-f
--      radix.36.u      Radix 36 of 0-9 A-Z         Also known as radix.36
--      radix.36.l      Radix 36 of 0-9 a-z
--      radix.62        Radix 62 of 0-9 A-Z a-z
--      radix.64        Radix 62 of 0-9 A-Z a-z @ _
--
-- To create new alphabets enter a new entry into id.meta with it's name, the
-- ival set to the radix and tval to the alphabet (must be ival characters long).
--

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('hex',
        16,
        '0123456789ABCDEF',
        'Hexadecimal upper case alphabet')
ON CONFLICT DO NOTHING;

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('hex.u',
        16,
        '0123456789ABCDEF',
        'Hexadecimal upper case alphabet')
ON CONFLICT DO NOTHING;

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('hex.l',
        16,
        '0123456789abcdef',
        'Hexadecimal lower case alphabet')
ON CONFLICT DO NOTHING;

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('radix.36',
        36,
        '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ',
        'Base 36 upper case alphabet')
ON CONFLICT DO NOTHING;

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('radix.36.u',
        36,
        '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ',
        'Base 36 upper case alphabet')
ON CONFLICT DO NOTHING;

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('radix.36.l',
        36,
        '0123456789abcdefghijklmnopqrstuvwxyz',
        'Base 36 lower case alphabet')
ON CONFLICT DO NOTHING;

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('radix.62',
        62,
        '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz',
        'Base 62 alphabet')
ON CONFLICT DO NOTHING;

-- 62m - alternat to 62 where case is mixed to give an alternate sequence of chars
INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('radix.62m',
        62,
        '0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz',
        'Base 62 alphabet variant')
ON CONFLICT DO NOTHING;

INSERT INTO id.meta (name, ival, tval, descr)
VALUES ('radix.64',
        64,
        '0123456789ABCDEFGHIJKLMNOPQRSTIVWXYZabcdefghijklmnopqrstuvwxyz@_',
        'Base 64 alphabet')
ON CONFLICT DO NOTHING;

--
-- digits       the BIGINT to encode. NULL will return NULL.
--
-- _radix       The alphabet to use, defaults to 'radix.36'
--
-- min_width    The minimum width of the output. Default 1.
--              If output is shorter then '0' is prefixed until the output is
--              of that length.
--
CREATE OR REPLACE FUNCTION id.encode(digits BIGINT, _radix NAME = 'radix.36', min_width INTEGER = 1)
    RETURNS TEXT AS
$$
DECLARE
    chars CHAR[];
    radix BIGINT;
    ret   VARCHAR := '';
    val   BIGINT  := digits;
BEGIN
    IF digits IS NULL THEN
        RETURN NULL;
    END IF;

    chars := regexp_split_to_array(id.meta_text(_radix), '');
    radix := id.meta_int(_radix);

    IF val < 0 THEN
        val := val * -1;
    END IF;

    WHILE val != 0
        LOOP
            ret := chars[(val % radix) + 1] || ret;
            val := val / radix;
        END LOOP;

    IF min_width > 0 AND char_length(ret) < min_width THEN
        ret := lpad(ret, min_width, chars[1]);
    END IF;

    RETURN ret;
END;
$$ LANGUAGE plpgsql
   IMMUTABLE;

-- ================================================================================
-- Decode an encoded string into a BIGINT
--
-- str      String to decode
-- _radix   The alphabet to use, default 'radix.36'
--
CREATE OR REPLACE FUNCTION id.decode(str TEXT, _radix NAME = 'radix.36')
    RETURNS BIGINT AS
$$
DECLARE
    chars TEXT    := id.meta_text(_radix);
    radix BIGINT  := id.meta_int(_radix);
    len   INTEGER := length(str);
    idx   INTEGER := 1;
    val   BIGINT  := 0;
BEGIN
    IF str IS NULL THEN
        RETURN NULL;
    END IF;

    WHILE idx <= len
        LOOP
            val := (val * radix) + strpos(chars, substr(str, idx, 1)) - 1;
            idx := idx + 1;
        END LOOP;

    RETURN val;
END;
$$ LANGUAGE plpgsql
   IMMUTABLE;


-- ================================================================================
-- Trigger that ensures id is unique and immutable once set

CREATE OR REPLACE FUNCTION id.id_trigger()
    RETURNS TRIGGER
AS
$$
DECLARE
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- CREATE new id/key pair
        NEW.id = id.instagram();
        --NEW.key = id.encode(NEW.id,'radix.62');
    ELSIF TG_OP = 'UPDATE' THEN
        -- Check immutable fields
        IF OLD.id != NEW.id THEN
            RAISE EXCEPTION 'id is immutable: %s', OLD.id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Example of what to use to add to a table
-- CREATE TRIGGER iot_chdata_idgen
--     BEFORE INSERT OR UPDATE
--     ON iot.chdata
--     FOR EACH ROW
--     EXECUTE PROCEDURE id.id_key_trigger();

-- ================================================================================
-- Trigger to manage created column in a table
--
-- To use simply add this before insert or update of a table
-- which must have a created column present
--
CREATE OR REPLACE FUNCTION id.created_trigger()
    RETURNS TRIGGER
AS
$$
DECLARE
BEGIN
    -- Ensure created is set but immutable
    IF TG_OP = 'INSERT' AND NEW.created IS NULL THEN
        -- Ensure created is set
        NEW.created = now();
    ELSIF TG_OP = 'UPDATE' AND OLD.created != NEW.created THEN
        RAISE EXCEPTION 'Created is immutable';
    ELSIF NEW.created IS NULL THEN
        RAISE EXCEPTION 'Created on %s cannot be null', TG_TABLE_NAME;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ================================================================================
-- Trigger to manage created/modified columns in a table
--
-- To use simply add this before insert or update of a table
-- which must have a created and modified columns present
--
CREATE OR REPLACE FUNCTION id.modified_trigger()
    RETURNS TRIGGER
AS
$$
DECLARE
BEGIN
    -- Ensure created is set but immutable
    IF TG_OP = 'INSERT' THEN

        IF NEW.created IS NULL THEN
            NEW.created = now();
        END IF;

        IF NEW.modified IS NULL THEN
            NEW.modified = now();
        END IF;

    ELSIF TG_OP = 'UPDATE' THEN

        -- Ensure created is immutable
        IF OLD.created != NEW.created THEN
            RAISE EXCEPTION 'Created is immutable';
        ELSIF NEW.created IS NULL THEN
            RAISE EXCEPTION 'Created on %s cannot be null', TG_TABLE_NAME;
        END IF;

        -- Ensure modified is set but change only if there really is a change
        -- This ensures modified is the timestamp of the last change to the row,
        IF NEW.modified IS NULL OR NEW != OLD THEN
            NEW.modified = now();
        END IF;

    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ================================================================================
CREATE OR REPLACE FUNCTION id.encode(pdt DATE)
    RETURNS TEXT AS
$$
SELECT id.encode((FLOOR(EXTRACT(EPOCH FROM pdt) * 1000) - id.meta_bigint('epoch'))::BIGINT, 'radix.62');
$$ LANGUAGE SQL;

-- ================================================================================
CREATE OR REPLACE FUNCTION id.encode(pdt TIME)
    RETURNS TEXT AS
$$
SELECT id.encode((FLOOR(EXTRACT(EPOCH FROM pdt) * 1000) - id.meta_bigint('epoch'))::BIGINT, 'radix.62');
$$ LANGUAGE SQL;

-- ================================================================================
CREATE OR REPLACE FUNCTION id.encode(pdt TIMESTAMP WITHOUT TIME ZONE)
    RETURNS TEXT AS
$$
SELECT id.encode((FLOOR(EXTRACT(EPOCH FROM pdt) * 1000) - id.meta_bigint('epoch'))::BIGINT, 'radix.62');
$$ LANGUAGE SQL;

-- ================================================================================
CREATE OR REPLACE FUNCTION id.encode(pdt TIMESTAMP WITH TIME ZONE)
    RETURNS TEXT AS
$$
SELECT id.encode((FLOOR(EXTRACT(EPOCH FROM pdt) * 1000) - id.meta_bigint('epoch'))::BIGINT, 'radix.62');
$$ LANGUAGE SQL;
