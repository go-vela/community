# v0.7.0 migration

When migrating from Vela version [v0.6.0](https://github.com/go-vela/community/blob/master/releases/v0.6.0.md) to [v0.7.0](https://github.com/go-vela/community/blob/master/releases/v0.7.0.md)  there is a breaking change for all native secrets. The change was to [encrypt all secret values](https://github.com/go-vela/community/issues/100) stored in the native engine. This means the endpoints in v0.7.0 where not compatible to read/write old secrets with the new encryption process. Secrets that reside in the database need to be encrypted.

Instead of building backwards compatibility we want all secrets to be encrypted which meant we need to supply a script that can be used during the version upgrade to assist admins with keeping their secrets create pre [v0.6.0](https://github.com/go-vela/community/blob/master/releases/v0.6.0.md) compatible with later versions. 

**Script Requirements:**

* [Go Installed](https://golang.org/doc/install)
  * Along with packages referenced in `main.go` file in this directory
  * Tip: packages can be downloaded with `go.mod`
* Address to vela server
* encryption key used for v0.7.0 installation
* Vela user account token with admin access

**Usage:**

```sh
# run the script within this directory
export VELA_ADDR=<server_addr>
export VELA_KEY=<server_key>
export VELA_TOKEN=<admin_token>

# execute program
go run main.go
```