#!/bin/bash

gcloud functions deploy \
  receive_measurement \
  --entry-point ReceiveMeasurement \
  --runtime go111 \
  --trigger-resource ruuvitag-measurements \
  --trigger-event google.pubsub.topic.publish \
  --set-env-vars DATASTORE_PROJECT_ID=ruuvitag-212713
