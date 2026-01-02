# Stage 1: Go builder
FROM golang:1-alpine3.22 AS go-builder

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev

WORKDIR /build
COPY . .
RUN ./build-go.sh

# Stage 2: Runtime
FROM alpine:3.22

ENV UID=1337 \
    GID=1337

RUN apk add --no-cache ffmpeg su-exec ca-certificates olm bash jq yq-go curl

COPY --from=go-builder /build/mautrix-twitter /usr/bin/mautrix-twitter
COPY --from=go-builder /build/docker-run.sh /docker-run.sh
VOLUME /data

CMD ["/docker-run.sh"]
