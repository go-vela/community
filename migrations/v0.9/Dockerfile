# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

################################################################################
## docker build --no-cache --target binary -t target/vela-migration:binary .  ##
################################################################################

FROM golang:latest@sha256:bc5f0b5e43282627279fe5262ae275fecb3d2eae3b33977a7fd200c7a760d6f1 as binary

COPY . /opt/go-vela/community

WORKDIR /opt/go-vela/community

RUN make clean

RUN make build-linux-static

################################################################################
##  docker build --no-cache --target binary -t target/vela-migration:certs .  ##
################################################################################

FROM alpine@sha256:82d1e9d7ed48a7523bdebc18cf6290bdb97b82302a8a9c27d4fe885949ea94d1 as certs

RUN apk add --update --no-cache ca-certificates

################################################################################
##  docker build --no-cache --target binary -t target/vela-migration:local .  ##
################################################################################

FROM scratch

COPY --from=binary /opt/go-vela/community/release/linux/amd64/vela-migration /bin/vela-migration
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/bin/vela-migration"] 