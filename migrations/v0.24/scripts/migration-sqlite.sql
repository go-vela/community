/*
    VELA MIGRATION v0.23.4 --> v0.24.0

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.24.X
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

-- Create dashboards table
CREATE TABLE
IF NOT EXISTS
dashboards (
	id            UUID PRIMARY KEY,
	name          TEXT,
	created_at    INTEGER,
	created_by    TEXT,
	updated_at    INTEGER,
	updated_by    TEXT,
	admins        TEXT,
	repos         TEXT
);

-- Add dashboards to users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS dashboards VARCHAR(5000)
;

-- Create settings table
CREATE TABLE
IF NOT EXISTS
settings (
	id                		SERIAL PRIMARY KEY,
	compiler          		TEXT,
	queue         	  		TEXT,
	repo_allowlist	  		TEXT,
	schedule_allowlist	    TEXT,
	created_at         		INTEGER,
	updated_at         		INTEGER,
	updated_by         		TEXT
);

-- Create jwks table
CREATE TABLE
IF NOT EXISTS
jwks (
	id     TEXT PRIMARY KEY,
	active BOOLEAN,
	key    TEXT
);

-- Add report_as to steps table
ALTER TABLE steps
    ADD COLUMN IF NOT EXISTS report_as TEXT
;

-- Add error to schedules table
ALTER TABLE schedules
    ADD COLUMN IF NOT EXISTS error TEXT
;

-- Remove allow_push, allow_pull, allow_tag, allow_comment, and allow_deploy from repos table
ALTER TABLE repos
    DROP COLUMN IF EXISTS allow_push,
    DROP COLUMN IF EXISTS allow_pull,
    DROP COLUMN IF EXISTS allow_tag,
    DROP COLUMN IF EXISTS allow_comment,
    DROP COLUMN IF EXISTS allow_deploy;

-- Remove events from secrets table
ALTER TABLE secrets
    DROP COLUMN IF EXISTS events;

-- Add sender_scm_id to builds table
ALTER TABLE builds
    ADD COLUMN IF NOT EXISTS sender_scm_id VARCHAR(250)
;

/*
___  _  _ _    
|  \ |\/| |    
|__/ |  | |___           
*/

-- Set sender_scm_id to '0' in builds table
UPDATE builds 
    SET sender_scm_id = '0'
;