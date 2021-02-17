# v0.7 migration

_Please note: this notice applies when upgrading to the latest v0.7.x release._

When migrating from Vela version [v0.6.0](https://github.com/go-vela/community/blob/master/releases/v0.6.0.md) to [v0.7.2](https://github.com/go-vela/community/blob/master/releases/v0.7.2.md) the Vela administrator will want to ensure the following actions are being performed:

1. Update tables in the Postgres database
   * `ALTER TABLE users ADD COLUMN IF NOT EXISTS refresh_token VARCHAR(500);`
   * `ALTER TABLE builds ADD COLUMN IF NOT EXISTS deploy_payload VARCHAR(2000);`
   * `ALTER TABLE workers ADD COLUMN IF NOT EXISTS build_limit INTEGER;`

1. Configure the Vela OAuth App's callback (in GitHub) to point to `<vela-server>/authenticate` (do not use the UI for the address)

1. If you have configured the `VELA_PORT` environment variable, ensure that it does not include a leading `:` or remove it altogether if you are fine with the `8080` default setting for the port value