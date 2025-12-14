#!/bin/sh
set -e

# Build Juicebox FFI library
./build-rust.sh

# Copy FFI artifacts to project root and pkg/juicebox
cp -f pkg/juicebox/juicebox-sdk/target/release/libjuicebox_sdk_ffi.a .
cp -f pkg/juicebox/juicebox-sdk/swift/Sources/JuiceboxSdkFfi/juicebox-sdk-ffi.h pkg/juicebox/

# Build Go with LIBRARY_PATH set
LIBRARY_PATH=.:$LIBRARY_PATH ./build-go.sh "$@"
