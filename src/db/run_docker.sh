#!/bin/bash

docker-compose down
docker volume prune
docker-compose build
docker-compose up --force-recreate
# ./run_docker.sh