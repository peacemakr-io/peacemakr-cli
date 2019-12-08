#!/bin/sh

set -ex

docker build . -f Dockerfile -t peacemakr-cli
docker build . -f Dockerfile-test -t peacemakr-cli-test

# Release artifacts to aws ecr
if [ "$1" == "-release" ]; then
  TAG="latest"
  if [[ ! -z "${2}" ]]; then
    TAG=${2}
  fi
  
  docker tag peacemakr-cli:latest peacemakr/peacemakr-cli:${TAG}

  echo "YOU ARE ABOUT TO DEPLOY TO PRODUCTION DOCKERHUB,"
  echo "please verify the the image is correct,"
  docker images | grep peacemakr/peacemakr-cli.*${TAG}
  echo ""
  read -p "(press enter to deploy)"
  set -x


  docker push peacemakr/peacemakr-cli:${TAG}
fi

