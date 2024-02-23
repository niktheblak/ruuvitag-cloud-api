#!/usr/bin/env bash

set -e

pushd .

cd ../..

gcloud functions deploy go-http-function \
  --gen2 \
  --runtime=go121 \
  --region=europe-north1 \
  --source=. \
  --env-vars-file scripts/gcp/.env \
  --entry-point=ReceiveMeasurement \
  --trigger-topic=projects/ruuvitag-415112/topics/measurements

popd
