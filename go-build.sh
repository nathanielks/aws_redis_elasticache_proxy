#!/usr/bin/env bash
GOOS="${GOOS:-linux}"
GOARCH="${GOARCH:-386}"
OUTFILE="${OUTFILE:-redis_auth_proxy}"
docker run --rm -v "$PWD":/usr/src/app -w /usr/src/app golang:1.8 \
  env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -v -o "${OUTFILE}"
