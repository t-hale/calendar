#!/bin/bash

./scripts/build_proto.sh

do
docker build -t "${_DOCKER_IMAGE}" .
docker push "${_DOCKER_IMAGE}"