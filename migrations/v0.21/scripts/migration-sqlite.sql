/*
    VELA MIGRATION v0.20 --> v0.21

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.21.X
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

-- Create build_executables table
CREATE TABLE
IF NOT EXISTS
build_executables (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    build_id         INTEGER,
    data             BLOB,
    UNIQUE(build_id)
);

-- Add branch to schedules table
ALTER TABLE schedules 
    ADD COLUMN IF NOT EXISTS branch TEXT
;

/*
___  _  _ _    
|  \ |\/| |    
|__/ |  | |___           
*/

-- Set branch to repo default branch for existing schedules
UPDATE schedules 
    SET branch = r.branch 
    FROM (SELECT id, branch FROM repos) r WHERE schedules.repo_id = r.id AND schedules.branch IS NULL
; 