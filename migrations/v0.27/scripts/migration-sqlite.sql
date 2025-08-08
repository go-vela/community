/*
    VELA MIGRATION v0.26.x --> v0.27.x

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.27.x
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

-- Start table alter transaction
BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS secret_repo_allowlists (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    secret_id       INTEGER,
    repo            VARCHAR(500),
    UNIQUE(secret_id, repo)
);

-- Add route to builds table
ALTER TABLE builds
    ADD COLUMN IF NOT EXISTS route VARCHAR(250)
;

-- Add created_at timeout to logs table
ALTER TABLE logs
    ADD COLUMN IF NOT EXISTS created_at INTEGER
;

-- Add custom_props to repos table
ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS custom_props TEXT
;

-- Add scm to settings table
ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS scm TEXT
;

-- Add max_dashboard_repos to settings table
ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS max_dashboard_repos INTEGER
;

-- Add queue_restart_limit to settings table
ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS queue_restart_limit INTEGER
;

-- Add secret_repo_allowlists_secret_id index
CREATE INDEX IF NOT EXISTS secret_repo_allowlists_secret_id ON secret_repo_allowlists (secret_id)
;

-- Add logs_created_at index
CREATE INDEX IF NOT EXISTS logs_created_at ON logs (created_at)
;

-- Add logs_service_id index
CREATE INDEX IF NOT EXISTS logs_service_id ON logs (service_id)
;

-- Add logs_step_id index
CREATE INDEX IF NOT EXISTS logs_step_id ON logs (step_id)
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

-- Backfill created_at in logs table
-- Attempt to backfill from steps and services first for accuracy
-- Fallback to builds (via step/service) if still NULL
-- Final fallback to now() for any remaining NULLs

-- Backfill from steps
UPDATE logs
SET created_at = (
    SELECT s.created
    FROM steps s
    WHERE s.id = logs.step_id
        AND s.created IS NOT NULL
)
WHERE created_at IS NULL
AND step_id IS NOT NULL
AND EXISTS (
    SELECT 1 FROM steps s
    WHERE s.id = logs.step_id
        AND s.created IS NOT NULL
	)
;

-- Backfill from services
UPDATE logs
SET created_at = (
    SELECT sv.created
    FROM services sv
    WHERE sv.id = logs.service_id
        AND sv.created IS NOT NULL
)
WHERE created_at IS NULL
AND service_id IS NOT NULL
AND EXISTS (
    SELECT 1 FROM services sv
    WHERE sv.id = logs.service_id
        AND sv.created IS NOT NULL
	)
;

-- Optional: backfill from builds (via step) if still NULL
UPDATE logs
SET created_at = (
    SELECT b.created
    FROM steps s
    JOIN builds b ON b.id = s.build_id
    WHERE s.id = logs.step_id
        AND b.created IS NOT NULL
    LIMIT 1
)
WHERE created_at IS NULL
AND step_id IS NOT NULL
AND EXISTS (
    SELECT 1
    FROM steps s
    JOIN builds b ON b.id = s.build_id
    WHERE s.id = logs.step_id
        AND b.created IS NOT NULL
)
;

-- Optional: backfill from builds (via service) if still NULL
UPDATE logs
SET created_at = (
    SELECT b.created
    FROM services sv
    JOIN builds b ON b.id = sv.build_id
    WHERE sv.id = logs.service_id
        AND b.created IS NOT NULL
    LIMIT 1
)
WHERE created_at IS NULL
AND service_id IS NOT NULL
AND EXISTS (
    SELECT 1
    FROM services sv
    JOIN builds b ON b.id = sv.build_id
    WHERE sv.id = logs.service_id
        AND b.created IS NOT NULL
)
;

-- Final fallback to now() for any remaining NULLs
UPDATE logs
SET created_at = CAST(strftime('%s','now') AS INTEGER)
WHERE created_at IS NULL
;

COMMIT;
