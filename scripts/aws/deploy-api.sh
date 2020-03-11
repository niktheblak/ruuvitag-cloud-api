#!/usr/bin/env bash

set -e

pushd .

cd ../..
GOOS=linux go build -o api cmd/aws/api/*.go
zip api.zip api

set +e

source scripts/aws/api.env

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

popd
