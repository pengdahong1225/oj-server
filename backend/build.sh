#!/bin/bash

all_services=("db-service" "judge-service" "question-service")
function build() {
    echo "build $1 start ......."
    go build -o ./app/$1/service ./app/$1/cmd
    echo "build $1 end ......."
}

if [[ $# -eq 0 ]]; then
    # 使用默认的服务列表
    services=("${all_services[@]}")
else
    # 使用传递的参数作为服务列表
    services=("$@")
fi

echo "go mod tidy"
go mod tidy

for srv in "${services[@]}"; do
    build ${srv}
done
