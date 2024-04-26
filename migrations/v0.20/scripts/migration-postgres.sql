/*
    VELA MIGRATION v0.19 --> v0.20

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.20.X
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

ALTER TABLE workers 
    ADD COLUMN IF NOT EXISTS status VARCHAR(50),
    ADD COLUMN IF NOT EXISTS last_status_update_at INTEGER,
    ADD COLUMN IF NOT EXISTS running_build_ids VARCHAR(500),
    ADD COLUMN IF NOT EXISTS last_build_started_at INTEGER,
    ADD COLUMN IF NOT EXISTS last_build_finished_at INTEGER
;