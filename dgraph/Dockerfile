# ------------------------------------------------------------------------------
# Build plugins
# ------------------------------------------------------------------------------

# building the plugins requires a similar build environment to the dgraph build
# environment and the Go version must match exactly. The go version can be found
# in the `release.sh` file listed as `GOVERSION`.
# https://github.com/dgraph-io/dgraph/blob/master/contrib/release.sh
FROM golang:1.17.3-bullseye as build-plugins
WORKDIR /build
COPY ./plugins/exacti ./exacti
RUN set -e \
    && mkdir -p plugins \
    && go build -buildmode=plugin -o ./plugins/exacti.so ./exacti/main.go

# ------------------------------------------------------------------------------
# Create DGraph image with plugins included
# ------------------------------------------------------------------------------

FROM dgraph/dgraph:v21.12.0
COPY --from=build-plugins /build/plugins /plugins
