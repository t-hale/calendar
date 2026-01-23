#!/bin/bash

mkdir -p gen
protoc --proto_path=proto --go_out=gen --go_opt=paths=source_relative $(find proto -iname "*.proto")