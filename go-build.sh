#!/usr/bin/env bash

docker run --rm -v "$PWD":/usr/src/app -w /usr/src/app golang:1.8 \
  env GOOS=linux GOARCH=386 go build -v -o redis_auth_proxy
