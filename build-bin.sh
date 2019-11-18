#!/bin/sh

set -ex

docker build . -f Dockerfile -t peacemakr-cli
docker build . -f Dockerfile-test -t peacemakr-cli-test
