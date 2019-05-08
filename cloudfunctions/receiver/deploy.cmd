@echo off

gcloud functions deploy ^
  receive_measurement ^
  --entry-point ReceiveMeasurement ^
  --runtime go111 ^
  --trigger-resource ruuvitag-measurements ^
  --trigger-event google.pubsub.topic.publish ^
  --env-vars-file .env.yaml