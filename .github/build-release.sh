#!/usr/bin/env bash

set -eux

mkdir -p dist package

extension=""
if [ "$GOOS" = "windows" ]; then
  extension=".exe"
fi

binary_path="package/${BINARY_NAME}${extension}"

go build -trimpath -ldflags="-s -w" -o "$binary_path" .

archive_base="${BINARY_NAME}_${TAG#v}_${GOOS}_${GOARCH}"

if [ "$ARCHIVE_FORMAT" = "zip" ]; then
  (
    cd package
    zip -9 "../dist/${archive_base}.zip" "${BINARY_NAME}${extension}"
  )
else
  tar -C package -czf "dist/${archive_base}.tar.gz" "${BINARY_NAME}${extension}"
fi

sha256sum "dist/${archive_base}.${ARCHIVE_FORMAT}" > "dist/${archive_base}.${ARCHIVE_FORMAT}.sha256"

rm -rf package