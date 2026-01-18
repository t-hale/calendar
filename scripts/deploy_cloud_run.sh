#!/bin/zsh

IMAGE="us-east1-docker.pkg.dev/api-project-119360632367/calendar/main"

docker build --platform linux/amd64 -t ${IMAGE:?} .
docker push ${IMAGE:?}
pushd terraform
terraform apply --auto-approve