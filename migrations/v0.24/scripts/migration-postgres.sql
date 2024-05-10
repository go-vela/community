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
	name          VARCHAR(250),
	created_at    INTEGER,
	created_by    VARCHAR(250),
	updated_at    INTEGER,
	updated_by    VARCHAR(250),
	admins        JSON DEFAULT NULL,
	repos         JSON DEFAULT NULL
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
	compiler          		JSON DEFAULT NULL,
	queue         	  		JSON DEFAULT NULL,
	repo_allowlist	  		VARCHAR(1000),
	schedule_allowlist	    VARCHAR(1000),
	created_at         		INTEGER,
	updated_at         		INTEGER,
	updated_by         		VARCHAR(250)
);

-- Add report_as to steps table
ALTER TABLE steps
    ADD COLUMN IF NOT EXISTS report_as VARCHAR(250)
;

-- Add error to schedules table
ALTER TABLE schedules
    ADD COLUMN IF NOT EXISTS error VARCHAR(250)
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