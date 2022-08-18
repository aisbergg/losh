# ------------------------------------------------------------------------------
# Build Web Assets
# ------------------------------------------------------------------------------

FROM node:18-alpine as build-assets
RUN apk add --no-cache --no-progress build-base automake autoconf libpng-dev
WORKDIR /build
COPY ./web/package.json /web/yarn.lock ./
RUN yarn install
COPY ./web/intf/http/assets ./intf/http/assets
COPY ./web/webpack.mix.js ./
RUN yarn run prod


# ------------------------------------------------------------------------------
# Build App
# ------------------------------------------------------------------------------

FROM golang:1.18-alpine as build-app
WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
COPY --from=build-assets /build/build/assets web/build/assets
RUN set -e \
    && mkdir -p build \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -a -ldflags '-s -w -extldflags "-static"' -o build/losh-web web/main.go \
    && mkdir /upload


# ------------------------------------------------------------------------------
# Assemble
# ------------------------------------------------------------------------------

FROM scratch
COPY --from=build-app /build/build/losh-web /losh-web
COPY --from=build-app --chown=999:0 /upload /
VOLUME ["/upload"]
EXPOSE 3000
USER 999
ENTRYPOINT ["/losh-web"]
