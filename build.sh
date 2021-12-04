#!/bin/bash

#aix/ppc64 android/386 android/amd64 android/arm android/arm64 darwin/amd64 darwin/arm64 dragonfly/amd64 freebsd/386 freebsd/amd64 freebsd/arm freebsd/arm64 illumos/amd64 ios/amd64 ios/arm64 js/wasm linux/386 linux/amd64 linux/arm linux/arm64 linux/mips linux/mips64 linux/mips64le linux/mipsle linux/ppc64 linux/ppc64le linux/riscv64 linux/s390x netbsd/386 netbsd/amd64 netbsd/arm netbsd/arm64 openbsd/386 openbsd/amd64 openbsd/arm openbsd/arm64 openbsd/mips64 plan9/386 plan9/amd64 plan9/arm solaris/amd64 windows/386 windows/amd64 windows/arm windows/arm64 

mkdir -p dist

supported_platforms=(linux/amd64 windows/amd64)
rel_ver="${1:-0.0.0}"
for platf in "${supported_platforms[@]}"; do
  _OS=`echo "$platf" | cut -d'/' -f1`
  _ARCH=`echo "$platf" | cut -d'/' -f2`
  if [ "$_ARCH" == "amd64" ]; then
    ARCH='x86_64'
  fi
  if [ "$_OS" == 'windows' ]; then
    EXT='.exe'
    OS=''
  elif [ "$_OS" == 'linux' ]; then
    EXT=''
    OS=linux-
  fi
  #echo "$OS $ARCH $platf"
  binary="fmt-md-text-$rel_ver-$OS$ARCH$EXT"
  echo "building: $binary"

  GOOS="$platf"
  go build -o "dist/$binary" .
done



