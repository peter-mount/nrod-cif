-- ======================================================================
-- Called at the beginning of a transaction to start a full CIF import
-- ======================================================================

-- Begin an import
-- returns true if the import should be performed, false if not
CREATE OR REPLACE FUNCTION timetable.beginimport(
  pFileMainframeIdentity   CHAR(20),
  pDateOfExtract           TIMESTAMP WITHOUT TIME ZONE,
  pCurrentFileReference    CHAR(7),
  pLastFileReference       CHAR(7),
  pUpdate                  BOOLEAN,
  pUserStartDate           TIMESTAMP WITH TIME ZONE,
  pUserEndDate             TIMESTAMP WITH TIME ZONE
)
RETURNS BOOLEAN AS $$
DECLARE
  rec RECORD;
BEGIN
  SELECT INTO rec c.*
    FROM timetable.cif c
    ORDER BY c.id DESC
    LIMIT 1;
  -- Last import is newer then ignore
  IF FOUND AND (
      -- Too old then ignore
      rec.DateOfExtract > pDateOfExtract
      -- same date then ignore if of the same type but if different allow
      -- a full import if
      -- same date then ignore if an update - A full import overrides an update
      OR (
        rec.DateOfExtract = pDateOfExtract
        AND ( pUpdate = rec.Update OR NOT pUpdate )
      )
    ) THEN
    RETURN FALSE;
  END IF;

  INSERT INTO timetable.cif
    (
      FileMainframeIdentity,
      DateOfExtract,
      CurrentFileReference,
      LastFileReference,
      Update,
      UserStartDate,
      UserEndDate,
      DateOfImport
    ) VALUES (
      pFileMainframeIdentity,
      pDateOfExtract,
      pCurrentFileReference,
      pLastFileReference,
      pUpdate,
      pUserStartDate,
      pUserEndDate,
      NOW()
    );

  -- Full import?
  IF NOT pUpdate THEN
    -- Removes all schedules
    DELETE FROM timetable.schedule;

    -- Catches any orphaned entries - this should not occur
    DELETE FROM timetable.station;

    -- Remove tiplocs
    DELETE FROM timetable.tiploc;

    -- Reset sequences
    ALTER SEQUENCE timetable.schedule_id_seq RESTART;
  END IF;

  RETURN TRUE;
END;
$$ LANGUAGE PLPGSQL;
