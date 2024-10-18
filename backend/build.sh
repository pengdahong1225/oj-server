# 编译
services=("./db-service" "./judge-service" "./question-service")

cd ./app || { echo "无法进入目录 app"; exit 1; }

# 遍历文件夹名
for s in "${services[@]}"; do
    path="$s"
    # 进入文件夹
    cd "$path" || { echo "无法进入目录 $path"; exit 1; }

    # 执行命令
    echo "$path : go mod tidy"
    go mod tidy

    echo "$path : go build -o service ./cmd"
    go build -o service ./cmd

    # 返回上级目录
    cd ..
done