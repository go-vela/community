### How it works today

#### Worker Routine

0. Vela Server is running (`VELA_SECRET` in environment). Worker is able to be run on a host.

1. Worker starts
    -> Environment: `VELA_SERVER_SECRET` (identical to server's `VELA_SECRET`)

2. Worker App Start Up
    -> Worker starts up, creates its server

3. Worker Check In
    -> Spins up Go Routine (parallel to step 4) for check in
    -> Loops on a specified interval (15m for Target), calling to POST / GET worker with the symmetric token
       as auth header

4. Worker Pulls Builds From Queue
    -> In a separate Go Routine (parallel to step 3), the worker will loop until it's reached its build limit and pull builds from Redis Queue
    -> Requests build token from server based on build information pulled from queue item
    -> Instantiates an executor for each build with its own SDK Client (auth header of build token)

#### Admin Process

1. At a basic level, admin starts the worker container with correct environment
    -> Admin must provide `VELA_SERVER_SECRET`, which should match the server env `VELA_SECRET`

1. Target Enterprise level, an admin uses Worker Ansible script to start container
    -> Script picks up `VELA_SERVER_SECRET` from Vault and puts it in the container environment
    -> Starts docker container

#### Not-so-happy Paths

- Worker Fails Check In
    -> No consequence to this other than stale data in the workers table
    -> Worker continues to pull builds from the queue and execute them

- Worker Panics
    -> Relaunch Go App




### How we want it to work

#### Server Changes
    A. New token types (Worker Registration, Worker Auth)
    B. Token manager able to mint / validate new token types with claims.Subject being the worker host name
    C. Permissions check on PUT / POST worker to verify token is of correct type
    D. Exchange auth token for new auth token in endpoints from C
    E. Add endpoint to validate a token _from_ server _to_ worker

#### Worker Routine

0. Vela Server is running. Worker is able to be run on a host.

1. Worker Starts
    -> Optional to include a registration token in the environment (in this case, step 3 will be skipped)

2. Worker App Start Up
    -> Worker starts up, creates its server
    -> Attaches middleware with context
        - Give server access to check in loop

3. Worker halts operation until its /register endpoint is hit

#### Option A: spin up worker container with registration 