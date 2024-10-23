#!/bin/bash

services=("question-service" "judge-service" "db-service")

for srv in "${services[@]}"; do
    docker stop ${srv}
    docker rm -f ${srv}
done

docker-compose -f docker-compose.service.yml up -d