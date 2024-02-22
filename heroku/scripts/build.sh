#!/usr/bin/env bash

set -e

pushd .
cd ../..
go build -o bin/heroku/server cmd/heroku/*.go
popd
