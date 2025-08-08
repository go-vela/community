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

CREATE TABLE
IF NOT EXISTS
secret_repo_allowlists (
	id                 BIGSERIAL PRIMARY KEY,
	secret_id          BIGINT,
	repo               VARCHAR(500),
	UNIQUE(secret_id, repo)
);

-- Add route to builds table
ALTER TABLE builds
    ADD COLUMN IF NOT EXISTS route VARCHAR(250)
;

-- Add created_at timeout to logs table
ALTER TABLE logs
    ADD COLUMN IF NOT EXISTS created_at BIGINT
;

-- Add custom_props to repos table
ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS custom_props JSON DEFAULT NULL
;

-- Add scm to settings table
ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS scm JSON DEFAULT NULL
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
UPDATE logs l
SET created_at = s.created
FROM steps s
WHERE l.step_id = s.id
  AND l.created_at IS NULL
  AND s.created IS NOT NULL
;

-- Backfill from services
UPDATE logs l
SET created_at = sv.created
FROM services sv
WHERE l.service_id = sv.id
  AND l.created_at IS NULL
  AND sv.created IS NOT NULL
;

-- Optional: backfill from builds (via step/service) if still NULL
UPDATE logs l
SET created_at = b.created
FROM steps s
JOIN builds b ON b.id = s.build_id
WHERE l.step_id = s.id
  AND l.created_at IS NULL
  AND b.created IS NOT NULL
;

UPDATE logs l
SET created_at = b.created
FROM services sv
JOIN builds b ON b.id = sv.build_id
WHERE l.service_id = sv.id
  AND l.created_at IS NULL
  AND b.created IS NOT NULL
;

-- Final fallback to now() for any remaining NULLs
-- If you don't care for accuracy, you can simplify to just this step
-- without the previous backfill steps. It should be much faster.
UPDATE logs
SET created_at = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE created_at IS NULL
;

COMMIT;