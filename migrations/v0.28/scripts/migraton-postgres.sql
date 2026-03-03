/*
    VELA MIGRATION v0.27.x --> v0.28.x

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.28.x
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

-- Start table alter transaction
BEGIN TRANSACTION;

-- Add artifacts to pipelines table
ALTER TABLE pipelines
    ADD COLUMN IF NOT EXISTS artifacts BOOLEAN
;

-- Add merge_queue_events to repos table
ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS merge_queue_events VARCHAR(500)
;

-- Add hook_counter to repos table
ALTER TABLE repos
    ADD COLUMN IF NOT EXISTS hook_counter BIGINT
;

-- Add enable_shared_secrets to settings table
ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS enable_shared_secrets BOOLEAN
;

-- Add enable_org_secrets to settings table
ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS enable_org_secrets BOOLEAN
;

-- Add enable_repo_secrets to settings table
ALTER TABLE settings
    ADD COLUMN IF NOT EXISTS enable_repo_secrets BOOLEAN
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

-- Update hook_counter in repos table
UPDATE repos r
SET hook_counter = COALESCE((
    SELECT COUNT(*)
    FROM hooks h
    WHERE h.repo_id = r.id
), 0);

-- Save changes
COMMIT;