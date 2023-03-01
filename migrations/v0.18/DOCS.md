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
cd $HOME/go-vela/community/migrations/v0.18
```

* Set the environment variables for the database driver and configuration string in your local terminal:

```sh
# set the driver for the Vela database
export VELA_DATABASE_DRIVER=<database driver from Vela server>

# set the address for the Vela database
export VELA_DATABASE_ADDR=<database address from Vela server>
```

## Start

> NOTE: Please review the [setup section](#setup) before moving forward.

This section covers the commands required to get the Vela application running locally.

* Navigate to the repository code:

```sh
# change into the project directory
cd $HOME/go-vela/community/migrations/v0.18
```

### CLI

This method of running the application uses the Golang binary built from the source code.

* Build the Golang binary targeting different operating systems and architectures:

```sh
# execute the `build` target with `make`
make build

# This command will output binaries to the following locations:
#
# * $HOME/go-vela/community/migrations/v0.18/release/darwin/amd64/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/linux/amd64/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/linux/arm64/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/linux/arm/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/windows/amd64/vela-migration
```

* Run the Golang binary for the specific operating system and architecture:

```sh
# run the Go binary for a Darwin (MacOS) operating system with amd64 architecture
release/darwin/amd64/vela-migration --action.all

# run the Go binary for a Linux operating system with amd64 architecture
release/linux/amd64/vela-migration --action.all

# run the Go binary for a Linux operating system with arm64 architecture
release/linux/arm64/vela-migration --action.all

# run the Go binary for a Linux operating system with arm architecture
release/linux/arm/vela-migration --action.all

# run the Go binary for a Windows operating system with amd64 architecture
release/windows/amd64/vela-migration --action.all
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
#   -e VELA_DATABASE_DRIVER \
#   -e VELA_DATABASE_ADDR \
#   target/vela-migration:local
```

## Usage

> NOTE: Please review the [start section](#start) before moving forward.

This utility supports invoking the following actions when migrating to `v0.18.x`:

* `alter.tables` - runs the action responsible for altering database tables

### Alter Tables

#### CLI

This method of running the application uses the Golang binary built from the source code.

* Build the Golang binary targeting different operating systems and architectures:

```sh
# execute the `build` target with `make`
make build

# This command will output binaries to the following locations:
#
# * $HOME/go-vela/community/migrations/v0.18/release/darwin/amd64/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/linux/amd64/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/linux/arm64/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/linux/arm/vela-migration
# * $HOME/go-vela/community/migrations/v0.18/release/windows/amd64/vela-migration
```

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
# execute the `run-alter` target with `make`
make run-alter

# This command is functionally equivalent to:
#
# docker run --rm \
#   -e VELA_ALTER_TABLES=true \
#   -e VELA_DATABASE_DRIVER \
#   -e VELA_DATABASE_ADDR \
#   target/vela-migration:local
```
