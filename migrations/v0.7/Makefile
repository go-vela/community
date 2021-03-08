# Copyright (c) 2021 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

# The `clean` target is intended to clean the workspace
# and prepare the local changes for submission.
#
# Usage: `make clean`
.PHONY: clean
clean: tidy vet fmt fix

# The `build` target is intended to compile
# the Go source code into a binary.
#
# Usage: `make build`
.PHONY: build
build: build-darwin build-linux build-windows

# The `build-static` target is intended to compile
# the Go source code into a statically linked binary.
#
# Usage: `make build-static`
.PHONY: build-static
build-static: build-darwin-static build-linux-static build-windows-static

# The `tidy` target is intended to clean up
# the Go module files (go.mod & go.sum).
#
# Usage: `make tidy`
.PHONY: tidy
tidy:
	@echo
	@echo "### Tidying Go module"
	@go mod tidy

# The `vet` target is intended to inspect the
# Go source code for potential issues.
#
# Usage: `make vet`
.PHONY: vet
vet:
	@echo
	@echo "### Vetting Go code"
	@go vet ./...

# The `fmt` target is intended to format the
# Go source code to meet the language standards.
#
# Usage: `make fmt`
.PHONY: fmt
fmt:
	@echo
	@echo "### Formatting Go Code"
	@go fmt ./...

# The `fix` target is intended to rewrite the
# Go source code using old APIs.
#
# Usage: `make fix`
.PHONY: fix
fix:
	@echo
	@echo "### Fixing Go Code"
	@go fix ./...

# The `build-darwin` target is intended to compile the
# Go source code into a Darwin compatible (MacOS) binary.
#
# Usage: `make build-darwin`
.PHONY: build-darwin
build-darwin:
	@echo
	@echo "### Building release/darwin/amd64/vela-migration binary"
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 \
		go build -a \
		-o release/darwin/amd64/vela-migration \
		github.com/go-vela/community/migrations/v0.7

# The `build-linux` target is intended to compile the
# Go source code into a Linux compatible binary.
#
# Usage: `make build-linux`
.PHONY: build-linux
build-linux:
	@echo
	@echo "### Building release/linux/amd64/vela-migration binary"
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 \
		go build -a \
		-o release/linux/amd64/vela-migration \
		github.com/go-vela/community/migrations/v0.7
	@echo
	@echo "### Building release/linux/arm64/vela-migration binary"
	GOOS=linux CGO_ENABLED=0 GOARCH=arm64 \
		go build -a \
		-o release/linux/arm64/vela-migration \
		github.com/go-vela/community/migrations/v0.7
	@echo
	@echo "### Building release/linux/arm/vela-migration binary"
	GOOS=linux CGO_ENABLED=0 GOARCH=arm \
		go build -a \
		-o release/linux/arm/vela-migration \
		github.com/go-vela/community/migrations/v0.7

# The `build-windows` target is intended to compile the
# Go source code into a Windows compatible binary.
#
# Usage: `make build-windows`
.PHONY: build-windows
build-windows:
	@echo
	@echo "### Building release/windows/amd64/vela-migration binary"
	GOOS=windows CGO_ENABLED=0 GOARCH=amd64 \
		go build -a \
		-o release/windows/amd64/vela-migration \
		github.com/go-vela/community/migrations/v0.7

# The `build-darwin-static` target is intended to compile the
# Go source code into a statically linked, Darwin compatible (MacOS) binary.
#
# Usage: `make build-darwin-static`
.PHONY: build-darwin-static
build-darwin-static:
	@echo
	@echo "### Building release/darwin/amd64/vela-migration binary"
	GOOS=darwin CGO_ENABLED=0 GOARCH=amd64 \
		go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/darwin/amd64/vela-migration \
		github.com/go-vela/community/migrations/v0.7

# The `build-linux-static` target is intended to compile the
# Go source code into a statically linked, Linux compatible binary.
#
# Usage: `make build-linux-static`
.PHONY: build-linux-static
build-linux-static:
	@echo
	@echo "### Building release/linux/amd64/vela-migration binary"
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 \
		go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/linux/amd64/vela-migration \
		github.com/go-vela/community/migrations/v0.7
	@echo
	@echo "### Building release/linux/arm64/vela-migration binary"
	GOOS=linux CGO_ENABLED=0 GOARCH=arm64 \
		go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/linux/arm64/vela-migration \
		github.com/go-vela/community/migrations/v0.7
	@echo
	@echo "### Building release/linux/arm/vela-migration binary"
	GOOS=linux CGO_ENABLED=0 GOARCH=arm \
		go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/linux/arm/vela-migration \
		github.com/go-vela/community/migrations/v0.7

# The `build-windows-static` target is intended to compile the
# Go source code into a statically linked, Windows compatible binary.
#
# Usage: `make build-windows-static`
.PHONY: build-windows-static
build-windows-static:
	@echo
	@echo "### Building release/windows/amd64/vela-migration binary"
	GOOS=windows CGO_ENABLED=0 GOARCH=amd64 \
		go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/windows/amd64/vela-migration \
		github.com/go-vela/community/migrations/v0.7

# The `build-static-ci` target is intended to compile
# the Go source code into a statically linked binary
# when used within a CI environment.
#
# Usage: `make build-static-ci`
.PHONY: build-static-ci
build-static-ci:
	@echo
	@echo "### Building CI static release/vela-migration binary"
	@go build -a \
		-ldflags '-s -w -extldflags "-static"' \
		-o release/vela-migration \
		github.com/go-vela/community/migrations/v0.7

# The `docker-build` target is intended to build
# the Docker image for the utility.
#
# Usage: `make docker-build`
.PHONY: docker-build
docker-build:
	@echo
	@echo "### Building target/vela-migration:local image"
	@docker build --no-cache -t target/vela-migration:local .

# The `docker-run` target is intended to execute
# the Docker image for the utility.
#
# Usage: `make docker-run`
.PHONY: docker-run
docker-run:
	@echo
	@echo "### Executing target/vela-migration:local image"
	@docker run --rm \
		-e VELA_ALL=true \
		-e VELA_BUILD_LIMIT \
		-e VELA_CONCURRENCY_LIMIT \
		-e VELA_SECRET_LIMIT \
		-e VELA_DATABASE_DRIVER \
		-e VELA_DATABASE_CONFIG \
		-e VELA_DATABASE_CONNECTION_OPEN \
		-e VELA_DATABASE_CONNECTION_IDLE \
		-e VELA_DATABASE_CONNECTION_LIFE \
		-e VELA_DATABASE_COMPRESSION_LEVEL \
		-e VELA_DATABASE_ENCRYPTION_KEY \
		target/vela-migration:local

# The `docker-run-alter` target is intended to execute
# the Docker image for the utility.
#
# Usage: `make docker-run-alter`
.PHONY: docker-run-alter
docker-run-alter:
	@echo
	@echo "### Executing target/vela-migration:local image"
	@docker run --rm \
		-e VELA_ALTER_TABLES=true \
		-e VELA_DATABASE_DRIVER \
		-e VELA_DATABASE_CONFIG \
		-e VELA_DATABASE_CONNECTION_OPEN \
		-e VELA_DATABASE_CONNECTION_IDLE \
		-e VELA_DATABASE_CONNECTION_LIFE \
		target/vela-migration:local

# The `docker-run-compress` target is intended to execute
# the Docker image for the utility.
#
# Usage: `make docker-run-compress`
.PHONY: docker-run-compress
docker-run-compress:
	@echo
	@echo "### Executing target/vela-migration:local image"
	@docker run --rm \
		-e VELA_COMPRESS_LOGS=true \
		-e VELA_BUILD_LIMIT \
		-e VELA_CONCURRENCY_LIMIT \
		-e VELA_DATABASE_DRIVER \
		-e VELA_DATABASE_CONFIG \
		-e VELA_DATABASE_CONNECTION_OPEN \
		-e VELA_DATABASE_CONNECTION_IDLE \
		-e VELA_DATABASE_CONNECTION_LIFE \
		-e VELA_DATABASE_COMPRESSION_LEVEL \
		target/vela-migration:local

# The `docker-run-drop` target is intended to execute
# the Docker image for the utility.
#
# Usage: `make docker-run-drop`
.PHONY: docker-run-drop
docker-run-drop:
	@echo
	@echo "### Executing target/vela-migration:local image"
	@docker run --rm \
		-e VELA_DROP_INDEXES=true \
		-e VELA_DATABASE_DRIVER \
		-e VELA_DATABASE_CONFIG \
		-e VELA_DATABASE_CONNECTION_OPEN \
		-e VELA_DATABASE_CONNECTION_IDLE \
		-e VELA_DATABASE_CONNECTION_LIFE \
		target/vela-migration:local

# The `docker-run-encrypt` target is intended to execute
# the Docker image for the utility.
#
# Usage: `make docker-run-encrypt`
.PHONY: docker-run-encrypt
docker-run-encrypt:
	@echo
	@echo "### Executing target/vela-migration:local image"
	@docker run --rm \
		-e VELA_ENCRYPT_SECRETS=true \
		-e VELA_SECRETS_LIMIT \
		-e VELA_CONCURRENCY_LIMIT \
		-e VELA_DATABASE_DRIVER \
		-e VELA_DATABASE_CONFIG \
		-e VELA_DATABASE_CONNECTION_OPEN \
		-e VELA_DATABASE_CONNECTION_IDLE \
		-e VELA_DATABASE_CONNECTION_LIFE \
		-e VELA_DATABASE_ENCRYPTION_KEY \
		target/vela-migration:local