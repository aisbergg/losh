# ------------------------------------------------------------------------------
# Build Web Assets
# ------------------------------------------------------------------------------

FROM node:18-alpine as build-static
WORKDIR /build
COPY ./web/package.json /web/yarn.lock ./
RUN yarn install
COPY ./web/resources ./resources
COPY ./web/webpack.mix.js ./
RUN yarn run prod


# ------------------------------------------------------------------------------
# Build App
# ------------------------------------------------------------------------------

FROM golang:1.18-alpine as build-app
WORKDIR /build
COPY ./go.mod ./go.sum ./
COPY ./lib/vendor/ ./lib/vendor/
RUN go mod download
COPY ./ ./
RUN set -e \
    && mkdir -p build \
    && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static"' -o build/losh-web web/main.go \
    && mkdir /upload


# ------------------------------------------------------------------------------
# Assemble
# ------------------------------------------------------------------------------

FROM scratch
COPY --from=build-static /build/build/assets/static /assets/static
COPY --from=build-app /build/web/resources/views /assets/templates
COPY --from=build-app /build/build/losh-web /losh-web
COPY --from=build-app --chown=999:0 /upload /
VOLUME ["/upload"]
EXPOSE 8080
EXPOSE 8443
USER 999
ENTRYPOINT ["/losh-web"]
