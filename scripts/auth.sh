#!/bin/bash

gcloud config set project ${PROJECT}
#gcloud auth application-default set-quota-project ${PROJECT}
#gcloud config set billing/quota_project ${BILLING_ACCOUNT_ID}
gcloud auth application-default login --scopes https://www.googleapis.com/auth/calendar,https://www.googleapis.com/auth/cloud-platform

