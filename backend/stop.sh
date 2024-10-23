#!/bin/bash

if [[ $# -eq 0 ]]; then
    echo "Usage: ./stop.sh <service_name>"
    exit 1
fi


docker kill -s 10 $1
docker-compose -f docker-compose.service.yml stop $1
echo "Service stopped."
docker rm -f $1
