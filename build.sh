#!/bin/sh

set -ex

main() {
  cleanup
  setTags

  buildImage
  testImage
}

cleanup() {
  \rm -f report.xml || true
}
buildImage() {
  docker build --pull=true --tag ${DOCKER_HUB_ACCOUNT}/${SERVICE_NAME}:${TAG} .
}

testImage() {
  docker run -i -P $DOCKER_HUB_ACCOUNT/$SERVICE_NAME:$TAG ./test.sh > report.xml

  CIDFILE=.cidfile.tmp
  \rm "${CIDFILE}" || true
  docker run -P -d --cidfile="${CIDFILE}" $DOCKER_HUB_ACCOUNT/$SERVICE_NAME:$TAG

  trap 'docker rm -f $(cat "${CIDFILE}"); rm -f "${CIDFILE}";' EXIT

  docker logs -f $(cat "${CIDFILE}") &

}

setTags() {
  SERVICE_NAME="$(basename "$(pwd)")"
  DEV_TAG="dev-latest"
  TAG="latest"
  DOCKER_HUB_ACCOUNT="local"
  echo "${TAG} development environment configuration..."
}

main
