#!/bin/sh
set -e

# Initialize/update submodule if needed
git submodule update --init pkg/juicebox/juicebox-sdk

# Build the Juicebox FFI crate from workspace root
cd pkg/juicebox/juicebox-sdk
RUSTFLAGS="-Ctarget-feature=-crt-static" RUSTC_WRAPPER="" cargo build -p juicebox_sdk_ffi --release
