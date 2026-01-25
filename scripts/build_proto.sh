#!/bin/bash

set -o errexit nounset pipefail

mkdir -p gen
docker run -v "$(pwd)":/calendar \
  jaegertracing/protobuf \
  --proto_path=calendar/proto \
  --go_out=calendar/gen \
  --go_opt=paths=source_relative \
  'calendar/proto/*.proto'