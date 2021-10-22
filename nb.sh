#!/usr/bin/bash

# TODO: Change terminology master/slave
# This creates a new slave that connects to the master

docker stop backend_$1 &>/dev/null
docker rm backend_$1 &>/dev/null
echo "Building nest binary..."
go build
echo "Done."
echo "Building docker image (can take a while)..."
docker build . -f Dockerfile.backend -t backend:latest >/dev/null
echo "Done." 
echo "Starting backend_$1"
docker run --name backend_$1 --restart=always -p 80:80 -dt backend:latest >/dev/null
echo "Done." 
docker inspect backend_$1 --format '{{.NetworkSettings.Networks.bridge.IPAddress}}'
