#!/usr/bin/env bash

set -e

source .env

GOOS=linux go build -o receiver ../../cmd/aws/receiver/main.go
zip receiver.zip receiver
aws lambda create-function \
  --function-name receiver \
  --runtime go1.x \
  --zip-file fileb://receiver.zip \
  --handler receiver \
  --role "$ROLE"
rm receiver receiver.zip
