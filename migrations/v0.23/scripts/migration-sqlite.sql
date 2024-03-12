/*
    VELA MIGRATION v0.22.2 --> v0.23.0

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.23.X
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

-- Create deployments table
CREATE TABLE
IF NOT EXISTS
deployments (
    id           SERIAL PRIMARY KEY,
    repo_id      INTEGER,
    number       INTEGER,	
    url     	   TEXT,
    "commit"     TEXT,
    ref          TEXT,
    task         TEXT,
    target       TEXT,
    description  TEXT,
    payload      TEXT,
    created_at   INTEGER,
    created_by   TEXT,
    builds       TEXT,
    UNIQUE(repo_id, number)
);

-- Add approve_build and allow_events to repos table
ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS approve_build TEXT
;

ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS allow_events INTEGER
;

-- Add approved_at, approved_by, and deploy_number to builds table
ALTER TABLE builds
    ADD COLUMN IF NOT EXISTS approved_at INTEGER
;

ALTER TABLE builds
    ADD COLUMN IF NOT EXISTS approved_by TEXT
;

ALTER TABLE builds 
    ADD COLUMN IF NOT EXISTS deploy_number INTEGER
;

-- Add allow_events to secrets table
ALTER TABLE secrets
    ADD COLUMN IF NOT EXISTS allow_events INTEGER
;

-- Add allow_substitution to secrets table (part of v0.23.2)
ALTER TABLE secrets
    ADD COLUMN IF NOT EXISTS allow_substitution BOOLEAN
;

/*
___  _  _ _    
|  \ |\/| |    
|__/ |  | |___           
*/

-- Set approve_build to 'fork-always' in repos table
UPDATE repos 
    SET approve_build = 'fork-always'
;

-- Map values of allow_<x> to allow_events integer mask in repos table
UPDATE repos
    SET allow_events =
        (CASE WHEN allow_push = 1 THEN 1 ELSE 0 END) |
        (CASE WHEN allow_pull = 1 THEN 4 | 16 | 1024 ELSE 0 END) |
        (CASE WHEN allow_tag = 1 THEN 2 ELSE 0 END) |
        (CASE WHEN allow_deploy = 1 THEN 8192 ELSE 0 END) |
        (CASE WHEN allow_comment = 1 THEN 16384 | 32768 ELSE 0 END)
;

-- Map slice of events to allow_events integer mask in secrets table
UPDATE secrets
    SET allow_events =
        (CASE WHEN events LIKE '%push%' THEN 1 ELSE 0 END) |
        (CASE WHEN events LIKE '%pull_request%' THEN 4 | 16 | 1024 ELSE 0 END) |
        (CASE WHEN events LIKE '%tag%' THEN 2 ELSE 0 END) |
        (CASE WHEN events LIKE '%deployment%' THEN 8192 ELSE 0 END) |
        (CASE WHEN events LIKE '%comment%' THEN 16384 | 32768 ELSE 0 END) |
        (CASE WHEN events LIKE '%schedule%' THEN 65536 ELSE 0 END)
;

-- Match the field for the new allow_substitution setting with the existing allow_command setting
UPDATE secrets SET allow_substitution = CASE WHEN allow_command = false THEN false ELSE true END;
