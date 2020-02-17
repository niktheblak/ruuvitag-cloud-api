#!/usr/bin/env bash

set -e

GOOS=linux go build -o receiver ../cmd/receiver/main.go
zip receiver.zip receiver

set +e

source receiver.env

if [[ -z "${FUNCTION_NAME}" ]]; then
  FUNCTION_NAME=receiver
fi

if [[ -z "${CREATE}" ]]; then
  aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://receiver.zip \
    --publish
else
  aws lambda create-function \
    --function-name $FUNCTION_NAME \
    --runtime go1.x \
    --zip-file fileb://receiver.zip \
    --handler receiver \
    --role $ROLE
fi

rm receiver receiver.zip
