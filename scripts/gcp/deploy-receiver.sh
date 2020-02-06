#!/usr/bin/env bash

gcloud functions deploy \
  receive_measurement \
  --source ../../cmd/gcp/receiver \
  --entry-point ReceiveMeasurement \
  --runtime go113 \
  --trigger-resource ruuvitag-measurements \
  --trigger-event google.pubsub.topic.publish \
  --env-vars-file .env.yaml
