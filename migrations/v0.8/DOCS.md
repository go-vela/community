# Documentation

This document intends to provide information on how to get the Vela migration utility running locally.

## Prerequisites

* [Docker](https://docs.docker.com/install/) - building block for local development
* [Golang](https://golang.org/dl/) - for source code and [dependency management](https://github.com/golang/go/wiki/Modules)
* [Make](https://www.gnu.org/software/make/) - start up local development

## Setup

> NOTE: Please review the [prerequisites section](#prerequisites) before moving forward.

* Clone this repository to your workstation:

```sh
# clone the project
git clone git@github.com:go-vela/community.git $HOME/go-vela/community
```

* Navigate to the repository code:

```sh
# change into the project directory
cd $HOME/go-vela/community/migrations/v0.8
```

* Set the environment variables for the database driver and configuration string in your local terminal:

```sh
# set the driver for the Vela database
export VELA_DATABASE_DRIVER=<database driver from Vela server>

# set the configuration string for the Vela database
export VELA_DATABASE_CONFIG=<database config from Vela server>
```

* (OPTIONAL) Set the environment variables for the other database configuration in your local terminal:

```sh
# set the total number of open connections for the database (default: 0 - no limit)
export VELA_DATABASE_CONNECTION_OPEN=<database connection open from Vela server>

# set the total number of idle connections for the database (default: 2)
export VELA_DATABASE_CONNECTION_IDLE=<database connection idle from Vela server>

# set the duration for the life of the database connections (default: 30m)
export VELA_DATABASE_CONNECTION_LIFE=<database connection life from Vela server>
```

## Start

> NOTE: Please review the [setup section](#setup) before moving forward.

This section covers the commands required to get the Vela application running locally.

* Navigate to the repository code:

```sh
# change into the project directory
cd $HOME/go-vela/community/migrations/v0.8
```

### CLI

This method of running the application uses the Golang binary built from the source code.

* Build the Golang binary targeting different operating systems and architectures:

```sh
# execute the `build` target with `make`
make build

# This command will output binaries to the following locations:
#
# * $HOME/go-vela/community/migrations/v0.8/release/darwin/amd64/vela-migration
# * $HOME/go-vela/community/migrations/v0.8/release/linux/amd64/vela-migration
# * $HOME/go-vela/community/migrations/v0.8/release/linux/arm64/vela-migration
# * $HOME/go-vela/community/migrations/v0.8/release/linux/arm/vela-migration
# * $HOME/go-vela/community/migrations/v0.8/release/windows/amd64/vela-migration
```

* Run the Golang binary for the specific operating system and architecture:

```sh
# run the Go binary for a Darwin (MacOS) operating system with amd64 architecture
release/darwin/amd64/vela-migration

# run the Go binary for a Linux operating system with amd64 architecture
release/linux/amd64/vela-migration

# run the Go binary for a Linux operating system with arm64 architecture
release/linux/arm64/vela-migration

# run the Go binary for a Linux operating system with arm architecture
release/linux/arm/vela-migration

# run the Go binary for a Windows operating system with amd64 architecture
release/windows/amd64/vela-migration
```

### Docker

This method of running the application uses a Docker container built from the `Dockerfile`.

* Build the Docker image:

```sh
# execute the `docker-build` target with `make`
make docker-build

# This command is functionally equivalent to:
#
# docker build --no-cache -t target/vela-migration:local .
```

* Run the Docker image

```sh
# execute the `run` target with `make`
make run

# This command is functionally equivalent to:
#
# docker run --rm \
#   -e VELA_ACTION_ALL=true \
#   -e VELA_CONCURRENCY_LIMIT \
#   -e VELA_DATABASE_DRIVER \
#   -e VELA_DATABASE_CONFIG \
#   -e VELA_DATABASE_CONNECTION_OPEN \
#   -e VELA_DATABASE_CONNECTION_IDLE \
#   -e VELA_DATABASE_CONNECTION_LIFE \
#   -e VELA_DATABASE_ENCRYPTION_KEY \
#   target/vela-migration:local
```

## Usage

> NOTE: Please review the [start section](#start) before moving forward.

This utility supports invoking the following actions when migrating to `v0.8.x`:

* `all` - run all supported actions configured in the migration utility
* `alter.tables` - runs the action responsible for altering database tables
* `sync.counter` - runs the action responsible for syncing repo counter values

### Alter Tables

#### CLI

This method of running the application uses the Golang binary built from the source code.

* Run the Golang binary for the specific operating system and architecture:

```sh
# run the Go binary for a Darwin (MacOS) operating system with amd64 architecture
release/darwin/amd64/vela-migration --alter.tables

# run the Go binary for a Linux operating system with amd64 architecture
release/linux/amd64/vela-migration --alter.tables

# run the Go binary for a Linux operating system with arm64 architecture
release/linux/arm64/vela-migration --alter.tables

# run the Go binary for a Linux operating system with arm architecture
release/linux/arm/vela-migration --alter.tables

# run the Go binary for a Windows operating system with amd64 architecture
release/windows/amd64/vela-migration --alter.tables
```

#### Docker

This method of running the application uses a Docker container built from the `Dockerfile`.

* Run the Docker image

```sh
# execute the `run-alter` target with `make`
make run-alter

# This command is functionally equivalent to:
#
# docker run --rm \
#   -e VELA_ALTER_TABLES=true \
#   -e VELA_DATABASE_DRIVER \
#   -e VELA_DATABASE_CONFIG \
#   -e VELA_DATABASE_CONNECTION_OPEN \
#   -e VELA_DATABASE_CONNECTION_IDLE \
#   -e VELA_DATABASE_CONNECTION_LIFE \
#   target/vela-migration:local
```

### Sync Counter

* Set the environment variables for the other database configuration in your local terminal:

```sh
# sets the limit of concurrent processes used to operate on the database (default: 4)
export VELA_CONCURRENCY_LIMIT=<range of 1 - runtime.GOMAXPROCS>
```

#### CLI

This method of running the application uses the Golang binary built from the source code.

* Run the Golang binary for the specific operating system and architecture:

```sh
# run the Go binary for a Darwin (MacOS) operating system with amd64 architecture
release/darwin/amd64/vela-migration --sync.counter

# run the Go binary for a Linux operating system with amd64 architecture
release/linux/amd64/vela-migration --sync.counter

# run the Go binary for a Linux operating system with arm64 architecture
release/linux/arm64/vela-migration --sync.counter

# run the Go binary for a Linux operating system with arm architecture
release/linux/arm/vela-migration --sync.counter

# run the Go binary for a Windows operating system with amd64 architecture
release/windows/amd64/vela-migration --sync.counter
```

#### Docker

This method of running the application uses a Docker container built from the `Dockerfile`.

* Run the Docker image

```sh
# execute the `run-sync` target with `make`
make run-sync

# This command is functionally equivalent to:
#
# docker run --rm \
#   -e VELA_SYNC_COUNTER=true \
#   -e VELA_CONCURRENCY_LIMIT \
#   -e VELA_DATABASE_DRIVER \
#   -e VELA_DATABASE_CONFIG \
#   -e VELA_DATABASE_CONNECTION_OPEN \
#   -e VELA_DATABASE_CONNECTION_IDLE \
#   -e VELA_DATABASE_CONNECTION_LIFE \
#   target/vela-migration:local
```