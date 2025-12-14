# Stage 1: Rust builder
FROM rust:1-alpine AS rust-builder

RUN apk add --no-cache git make musl-dev

WORKDIR /build
COPY pkg/juicebox/juicebox-sdk/. pkg/juicebox/juicebox-sdk/.
COPY build-rust.sh .
RUN ./build-rust.sh

# Stage 2: Go builder
FROM golang:1-alpine3.22 AS go-builder

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev

WORKDIR /build
COPY . .
COPY --from=rust-builder /build/pkg/juicebox/juicebox-sdk/target/release/libjuicebox_sdk_ffi.a ./
COPY --from=rust-builder /build/pkg/juicebox/juicebox-sdk/swift/Sources/JuiceboxSdkFfi/juicebox-sdk-ffi.h pkg/juicebox/
ENV LIBRARY_PATH=.
RUN ./build-go.sh

# Stage 3: Runtime
FROM alpine:3.22

ENV UID=1337 \
    GID=1337

RUN apk add --no-cache ffmpeg su-exec ca-certificates olm bash jq yq-go curl

COPY --from=go-builder /build/mautrix-twitter /usr/bin/mautrix-twitter
COPY --from=go-builder /build/docker-run.sh /docker-run.sh
VOLUME /data

CMD ["/docker-run.sh"]
