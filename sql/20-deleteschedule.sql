-- ======================================================================
-- Trigger on schedule delete remove all entries from the other tables
-- ======================================================================
CREATE OR REPLACE FUNCTION timetable.scheddeleted()
RETURNS TRIGGER AS $$
BEGIN
  -- Delete the json
  DELETE FROM timetable.schedule_json sj
    WHERE sj.id = OLD.id;
  -- Delete the index
  DELETE FROM timetable.station s
    WHERE s.sid = OLD.id;

  RETURN OLD;
END;
$$ LANGUAGE PLPGSQL;

CREATE TRIGGER scheddeleted
  BEFORE DELETE ON timetable.schedule
  FOR EACH ROW
  EXECUTE PROCEDURE timetable.scheddeleted();
