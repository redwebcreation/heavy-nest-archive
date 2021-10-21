#!/usr/bin/bash

# TODO: Change terminology master/slave
# This creates a new slave that connects to the master

go build
cp ./nest ./backend/nest
cd backend
docker build . --tag backend:latest 
docker run -t backend:latest 