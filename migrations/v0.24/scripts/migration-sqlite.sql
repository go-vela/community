/*
    VELA MIGRATION v0.23.x --> v0.24.x

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.24.x
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

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
