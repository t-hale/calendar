#!/bin/bash

gcloud infra-manager deployments apply $_DEPLOYMENT \
  --service-account=$_SERVICE_ACCOUNT \
  --git-source-repo=$_GITHUB_REPO \
  --git-source-directory=terraform \
  --git-source-ref=master