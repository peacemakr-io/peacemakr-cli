#!/bin/sh

set -ex

docker build . -f Dockerfile-dependencies -t peacemakr-cli-dependencies
