#!/bin/bash

#aix/ppc64 android/386 android/amd64 android/arm android/arm64 darwin/amd64 darwin/arm64 dragonfly/amd64 freebsd/386 freebsd/amd64 freebsd/arm freebsd/arm64 illumos/amd64 ios/amd64 ios/arm64 js/wasm linux/386 linux/amd64 linux/arm linux/arm64 linux/mips linux/mips64 linux/mips64le linux/mipsle linux/ppc64 linux/ppc64le linux/riscv64 linux/s390x netbsd/386 netbsd/amd64 netbsd/arm netbsd/arm64 openbsd/386 openbsd/amd64 openbsd/arm openbsd/arm64 openbsd/mips64 plan9/386 plan9/amd64 plan9/arm solaris/amd64 windows/386 windows/amd64 windows/arm windows/arm64 

BRANCH=$(git branch --show-current)
GIT_HASH=$(git rev-parse --short HEAD)
GIT_HASH_LONG=$(git rev-parse HEAD)
BUILD_DATE=$(date -I)
BRANCH_TYPE=${BRANCH%%/*} #master,release
VERSION=

debug(){
  echo "$BRANCH"
  echo ""
  echo "$GIT_HASH"
  echo "$GIT_HASH_LONG"
  echo "$BUILD_DATE"
  echo "$VERSION"
  echo ""
}

usage(){
  echo "build.sh get_version"
}

get_version(){
  if [ "$BRANCH" == "master" ]; then
    BRANCH_TYPE=dev
    VERSION="$GIT_HASH-dev"
  elif [ "${BRANCH%%/*}" == 'release' ]; then
    VERSION="${BRANCH##*/}"
  elif [ "${BRANCH%%/*}" == 'rc' ]; then
    VERSION="${BRANCH##*/}"
  else
    return 1
  fi
  printf "%s" "$VERSION"
  return 0
}

case "$1" in
  get_version) get_version;;
  *) echo "invalid option '$1'" && usage && exit 1
esac
exit 0
