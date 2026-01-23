#!/bin/bash

./scripts/build_proto.sh

_DOCKER_IMAGE="${1:-us-east1-docker.pkg.dev/api-project-119360632367/calendar/main}"
docker build -t "${_DOCKER_IMAGE}" .
