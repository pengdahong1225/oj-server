#!/bin/bash

services=("db-service" "judge-service" "question-service")
function build() {
    echo "build $1 start ......."
    go build -o service ./app/cmd/$1
    echo "build $1 end ......."
}

echo "go mod tidy"
go mod tidy

for srv in "${services[@]}"; do
    build ${srv}
done