#!/bin/bash

./scripts/build_proto.sh

docker build -t "${_DOCKER_IMAGE}" .
docker push "${_DOCKER_IMAGE}"