/*
    VELA MIGRATION v0.18 --> v0.19

    Please note that this SQL file must be executed prior to upgrading Vela to version 0.19.X
*/

/*
___  ___  _    
|  \ |  \ |    
|__/ |__/ |___ 
*/

ALTER TABLE repos 
    ADD COLUMN IF NOT EXISTS topics TEXT;