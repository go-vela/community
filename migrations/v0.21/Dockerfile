# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

################################################################################
## docker build --no-cache --target binary -t target/vela-migration:binary .  ##
################################################################################

FROM golang:latest as binary

COPY . /opt/go-vela/community

WORKDIR /opt/go-vela/community

RUN make clean

RUN make build-linux-static

###############################################################################
##  docker build --no-cache --target certs -t target/vela-migration:certs .  ##
###############################################################################

FROM alpine:latest as certs

RUN apk add --update --no-cache ca-certificates

################################################################
##  docker build --no-cache -t target/vela-migration:local .  ##
################################################################

FROM scratch

COPY --from=binary /opt/go-vela/community/release/linux/amd64/vela-migration /bin/vela-migration
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/bin/vela-migration"] 