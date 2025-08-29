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

-- Save changes
COMMIT;


/*
___  _  _ _    
|  \ |\/| |    
|__/ |  | |___           
*/

-- Start transaction
BEGIN TRANSACTION;

-- Set default scm role mapping in settings table
UPDATE settings
    SET scm = '{"repo_role_map":{"admin":"admin","maintain":"write","read":"read","triage":"read","write":"write"},"org_role_map":{"admin":"admin","member":"read"},"team_role_map":{"maintainer":"admin","member":"read"}}'
    WHERE scm IS NULL
;

-- Set max_dashboard_repos to 10 (or configured default max_dashboard_repos) in settings table
UPDATE settings
    SET max_dashboard_repos = 10
    WHERE max_dashboard_repos IS NULL
;

-- Set queue_restart_limit to 30 (or configured default queue_restart_limit) in settings table
UPDATE settings
    SET queue_restart_limit = 30
    WHERE queue_restart_limit IS NULL
;

-- Backfill created_at in logs table
-- Attempt to backfill from steps and services first for accuracy
-- Fallback to builds (via step/service) if still NULL
-- Final fallback to clock_timestamp() for any remaining NULLs
-- IMPORTANT: for large logs tables this may take a while to complete;
-- consider running this in smaller batches if you run into locks/timeouts
UPDATE logs l
SET created_at = COALESCE(
        s.created,
        sv.created,
        b_from_step.created,
        b_from_service.created,
        EXTRACT(EPOCH FROM clock_timestamp())::BIGINT
    )
FROM logs l2
LEFT JOIN steps    s  ON l2.step_id = s.id
LEFT JOIN services sv ON l2.service_id = sv.id
LEFT JOIN builds b_from_step    ON b_from_step.id = s.build_id
LEFT JOIN builds b_from_service ON b_from_service.id = sv.build_id
WHERE l.id = l2.id
  AND l.created_at IS NULL;

COMMIT;

-- Add logs_created_at index AFTER backfilling created_at
-- Note: remove CONCURRENTLY if your logs table is partitioned
CREATE INDEX CONCURRENTLY IF NOT EXISTS logs_created_at ON logs (created_at)
;

-- Analyze logs table to update statistics
ANALYZE logs
;