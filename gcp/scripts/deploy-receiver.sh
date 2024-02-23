#!/usr/bin/env bash

set -e

pushd .

cd ../receiver

gcloud functions deploy go-http-function \
  --gen2 \
  --runtime=go121 \
  --region=europe-north1 \
  --source=. \
  --env-vars-file ../scripts/env.json \
  --entry-point=ReceiveMeasurement \
  --trigger-topic=measurements

popd
