#!/bin/bash

if [[ $# -eq 0 ]]; then
    echo "Usage: ./stop.sh <service_name>"
    exit 1
fi

srv="backend_$1"

docker kill -s 10 $1
docker-compose -f docker-compose.service.yml stop $srv
echo "Service stopped."
docker rm -f $1