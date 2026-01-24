#!/bin/bash

printenv

gcloud auth list
gcloud config set project $PROJECT_NUMBER
gcloud auth login
gcloud container images describe $_DOCKER_IMAGE --format="value(image_summary.digest)" > image.txt
gcloud infra-manager deployments apply $_DEPLOYMENT \
  --service-account=$_SERVICE_ACCOUNT \
  --git-source-repo=$_GITHUB_REPO \
  --git-source-directory=terraform \
  --git-source-ref=master \
  --variables="image_tag=$(cat image.txt)"