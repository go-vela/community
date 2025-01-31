/*
    VELA MIGRATION v0.25.x --> v0.26.x

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.26.x
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

-- Start table alter transaction
BEGIN TRANSACTION;

-- Add warnings to pipelines table
ALTER TABLE pipelines
    ADD COLUMN IF NOT EXISTS warnings TEXT
;

-- Add approval timeout to repos table
ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS approval_timeout INTEGER
;

-- Add install_id to repos table
ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS install_id INTEGER
;

-- Add fork to builds table
ALTER TABLE builds
    ADD COLUMN IF NOT EXISTS fork BOOLEAN
;

-- Update hook error field to be larger
ALTER TABLE hooks
    ALTER COLUMN error TYPE VARCHAR(5000)
;

-- Delete builds_source index
DROP INDEX IF EXISTS builds_source
;

-- Add builds_event index
CREATE INDEX IF NOT EXISTS builds_event ON builds (event)
;

-- Add builds_repo_id_created index
CREATE INDEX IF NOT EXISTS builds_repo_id_created ON builds (repo_id, created)
;

-- Save changes
COMMIT;

/*
___  _  _ _    
|  \ |\/| |    
|__/ |  | |___           
*/

-- Start transaction
BEGIN TRANSACTION;

-- Set approval_timeout to 7 (or configured default approval_timeout) in repos table for all repos
UPDATE repos 
    SET approval_timeout = 7
;

COMMIT;