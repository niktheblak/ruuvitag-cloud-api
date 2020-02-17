#!/usr/bin/env bash

set -e

GOOS=linux go build -o api ../cmd/api/*.go
zip api.zip api

set +e

source api.env

if [[ -z "${FUNCTION_NAME}" ]]; then
  FUNCTION_NAME=api
fi

if [[ -z "${CREATE}" ]]; then
  aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://api.zip \
    --publish
else
  aws lambda create-function \
    --function-name $FUNCTION_NAME \
    --runtime go1.x \
    --zip-file fileb://api.zip \
    --handler api \
    --role $ROLE
fi

rm api api.zip
