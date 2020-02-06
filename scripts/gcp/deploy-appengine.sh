#!/usr/bin/env bash

gcloud app deploy ../../configs/gcp/app.yaml \
  --ignore-file .gcloudignore
