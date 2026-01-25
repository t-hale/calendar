#!/bin/bash

set -o errexit nounset pipefail

./scripts/build_proto.sh

docker build -t "${_DOCKER_IMAGE}" .
docker push "${_DOCKER_IMAGE}"