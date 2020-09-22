#!/bin/bash

# docker-compose down
docker container ls -a | grep "$(basename $PWD)\|go_test_" | awk '{print $1}' | xargs docker container rm -f
docker image ls -a | grep "$(basename $PWD)\|go_test_" | awk '{print $3}' | xargs docker image rm -f
docker network ls | grep "$(basename $PWD)\|go_test_" | awk '{print $1}' | xargs docker network rm
docker image prune