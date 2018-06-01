-- ==================================================
-- Create the database for nrod-cif
-- ==================================================

-- Core schema
\ir schema.sql

-- tiplocs
\ir tiploc.sql

-- schedules
\ir addschedule.sql
\ir deleteschedule.sql
\ir schedules.sql

-- manages CIF imports
\ir import.sql
