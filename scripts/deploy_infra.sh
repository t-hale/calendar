#!/bin/bash

usage() {
  echo "Usage: $0 --deployment <deployment_name> --service_account <email> --github_repo <repo_url> --git_source_directory <dir> --git_source_ref <branch/ref>"
  exit 1
}

# Parse named arguments
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --deployment) _DEPLOYMENT="$2"; shift ;;
    --service_account) _SERVICE_ACCOUNT="$2"; shift ;;
    --github_repo) _GITHUB_REPO="$2"; shift ;;
    --git_source_directory) _GIT_SOURCE_DIRECTORY="$2"; shift ;;
    --git_source_ref) _GIT_SOURCE_REF="$2"; shift ;;
    *) echo "Unknown parameter: $1"; usage ;;
  esac
  shift
done

# Fail if any parameter is missing
if [[ -z "$_DEPLOYMENT" ]] || [[ -z "$_SERVICE_ACCOUNT" ]] || [[ -z "$_GITHUB_REPO" ]] || [[ -z "$_GIT_SOURCE_DIRECTORY" ]] || [[ -z "$_GIT_SOURCE_REF" ]]; then
  echo "Error: Missing required arguments."
  usage
fi

gcloud infra-manager deployments apply "$_DEPLOYMENT" \
  --service-account="$_SERVICE_ACCOUNT" \
  --git-source-repo="$_GITHUB_REPO" \
  --git-source-directory="$_GIT_SOURCE_DIRECTORY" \
  --git-source-ref="$_GIT_SOURCE_REF"