#!/bin/bash
randomname=$(uuidgen)
container=wish/fmt-md-text

docker build -t "$container" . || exit 1
docker container create --name "$randomname" "$container" || exit 1
docker container cp "$randomname":/output/fmt-md-text ./ || exit 1
docker container rm "$randomname" || exit 1
