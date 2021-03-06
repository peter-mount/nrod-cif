-- ======================================================================
-- SQL Schema for the new timetable microservice
--
-- ======================================================================

DROP SCHEMA IF EXISTS timetable CASCADE;

CREATE SCHEMA IF NOT EXISTS timetable;

-- ======================================================================
-- record of db imports
-- ======================================================================
CREATE TABLE timetable.cif
(
    id                    SERIAL NOT NULL,
    FileMainframeIdentity CHAR(20),
    DateOfExtract         TIMESTAMP WITHOUT TIME ZONE,
    CurrentFileReference  CHAR(7),
    LastFileReference     CHAR(7),
    Update                BOOLEAN,
    UserStartDate         TIMESTAMP WITH TIME ZONE,
    UserEndDate           TIMESTAMP WITH TIME ZONE,
    -- when cif was imported
    DateOfImport          TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (id)
);
CREATE INDEX cif_c ON timetable.cif (FileMainframeIdentity);
CREATE INDEX cif_d ON timetable.cif (DateOfExtract);
CREATE INDEX cif_cd ON timetable.cif (FileMainframeIdentity, DateOfExtract);

-- ======================================================================
-- Tiploc location
--
-- The primary key is on id which is generated from the tiploc field by
-- the gettiplocid function
-- ======================================================================
CREATE TABLE timetable.tiploc
(
    -- Generated id based on tiploc by the gettiplocid function
    id          BIGINT                   NOT NULL,
    -- The TIPLOC id
    tiploc      VARCHAR(7)               NOT NULL,
    -- The CRS/3Alpha code or null
    crs         CHAR(3),
    -- The location's stanox id
    stanox      INTEGER,
    -- The full name of this tiploc
    name        VARCHAR(26),
    --- The NLC data for this location
    nlc         INTEGER,
    nlccheck    CHAR,
    nlcdesc     VARCHAR(16),
    -- true if this entry represents a station
    station     BOOLEAN                  NOT NULL DEFAULT FALSE,
    -- true if it's been "deleted". Entry kept for integrity purposes
    deleted     BOOLEAN                  NOT NULL DEFAULT FALSE,
    -- When this entry was created/modified
    dateextract TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX tiploc_tiploc
    ON timetable.tiploc (tiploc);

CREATE INDEX tiploc_crs
    ON timetable.tiploc (crs)
    WHERE crs IS NOT NULL;

CREATE INDEX tiploc_stanox
    ON timetable.tiploc (stanox)
    WHERE stanox IS NOT NULL;

CREATE INDEX tiploc_name
    ON timetable.tiploc (name)
    WHERE name IS NOT NULL;

-- Used for clustering related tiploc's together
CREATE UNIQUE INDEX tiploc_cluster
    ON timetable.tiploc (stanox, tiploc);

-- ======================================================================
-- schedule contains the searchable details of a schedule
-- ======================================================================
CREATE TABLE timetable.schedule
(
    -- generated id based on the primary key and the scheduleid function
    id        BIGINT   NOT NULL,
    -- Primary key for all schedules
    uid       CHAR(6)  NOT NULL,
    stp       CHAR     NOT NULL,
    startdate DATE     NOT NULL,
    -- end date so we can search by date range
    enddate   DATE     NOT NULL,
    -- The days of the week it's valid for
    dow       SMALLINT NOT NULL,
    -- entry date so we can optimise updates
    entrydate DATE     NOT NULL,
    PRIMARY KEY (uid, stp, startdate)
);
CREATE UNIQUE INDEX schedule_id ON timetable.schedule (id);
CREATE INDEX schedule_uid ON timetable.schedule (uid);
CREATE INDEX schedule_sd ON timetable.schedule (startdate);
CREATE INDEX schedule_ed ON timetable.schedule (enddate);
CREATE INDEX schedule_sed ON timetable.schedule (startdate, enddate);
CREATE INDEX schedule_used ON timetable.schedule (uid, stp, startdate, enddate);

-- ======================================================================
-- the schedule json
-- ======================================================================
CREATE TABLE timetable.schedule_json
(
    -- generated id, linked to the schedule table
    id       BIGINT NOT NULL REFERENCES timetable.schedule (id),
    -- The full json of this schedule
    schedule JSON   NOT NULL,
    PRIMARY KEY (id)
);

-- ======================================================================
-- Link between schedules and each individual station
--
-- When a schedule is added it's scanned and all stops are included in
-- this table.
--
-- When a search is performed for a station, this table is used and we
-- can then filter down the required schedules easily & quickly.
-- ======================================================================
CREATE TABLE timetable.station
(
    -- generated schedule id, linked to the schedule table
    sid       BIGINT   NOT NULL REFERENCES timetable.schedule (id),
    -- The position of this entry within the schedule
    ord       SMALLINT NOT NULL,
    -- The tiploc of this entry
    tid       BIGINT   NOT NULL,
    -- The schedule stp
    stp       CHAR     NOT NULL,
    -- used in searches, the date range for this entry
    startdate DATE     NOT NULL,
    enddate   DATE     NOT NULL,
    -- The days of the week it's valid for
    dow       SMALLINT NOT NULL,
    -- The time of the day at this point
    time      TIME     NOT NULL,
    PRIMARY KEY (sid, ord, tid)
);

CREATE INDEX station_i ON timetable.station (sid);
CREATE INDEX station_t ON timetable.station (tid);
CREATE INDEX station_td ON timetable.station (tid, startdate, enddate);
CREATE INDEX station_tdt ON timetable.station (tid, startdate, enddate, time);
CREATE INDEX station_tt ON timetable.station (tid, time);
--CREATE INDEX station_io ON timetable.station(scheduleId,ord);

-- ======================================================================
-- Schedule associations
-- ======================================================================
CREATE TABLE timetable.assoc
(
    id          SERIAL   NOT NULL,
    mainuid     CHAR(6)  NOT NULL,
    assocuid    CHAR(6)  NOT NULL,
    stp         CHAR     NOT NULL,
    startdate   DATE     NOT NULL,
    enddate     DATE     NOT NULL,
    dow         SMALLINT NOT NULL,
    cat         CHAR(2)  NOT NULL,
    dateInd     CHAR     NOT NULL,
    tid         BIGINT   NOT NULL,
    baseSuffix  CHAR     NOT NULL,
    assocSuffix CHAR     NOT NULL,
    assocType   CHAR     NOT NULL,
    -- entry date so we can optimise updates
    entrydate   DATE     NOT NULL,
    PRIMARY KEY (mainuid, assocuid, startdate, stp)
);

CREATE INDEX assoc_m ON timetable.assoc (mainuid);
CREATE INDEX assoc_a ON timetable.assoc (assocuid);
CREATE INDEX assoc_ma ON timetable.assoc (mainuid, assocuid);

CREATE INDEX assoc_mss ON timetable.assoc (mainuid, startdate, stp);
CREATE INDEX assoc_mses ON timetable.assoc (mainuid, startdate, enddate, stp);

CREATE INDEX assoc_ass ON timetable.assoc (assocuid, startdate, stp);
CREATE INDEX assoc_ases ON timetable.assoc (assocuid, startdate, enddate, stp);

CREATE INDEX assoc_cluster ON timetable.assoc (mainuid, assocuid, stp);
