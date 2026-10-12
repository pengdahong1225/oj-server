#!/bin/sh

set -e

cd "$(dirname "$0")"

COMPOSE_FILE=docker-compose.service.yaml

# ---- 选择 docker compose 命令：v2 插件优先，兼容 v1 独立二进制 ----
if docker compose version >/dev/null 2>&1; then
  DC="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  DC="docker-compose"
else
  echo "✘ 未检测到 docker compose，请先安装 docker 及 compose 插件" >&2
  exit 1
fi

# ---- 前置检查：网络 oj 必须由 deploy_env.sh 创建 ----
if ! docker network inspect oj >/dev/null 2>&1; then
  echo "✘ 未找到 docker 网络 oj，请先执行 ./deploy_env.sh" >&2
  exit 1
fi

# ---- 前置检查：二进制必须由 build.sh 编译 ----
missing=0
for name in gateway user judge problem; do
  if [ ! -f "./svr/$name/bin/$name" ]; then
    echo "✘ 缺少 ./svr/$name/bin/$name，请先执行 ./build.sh" >&2
    missing=1
  fi
done
if [ "$missing" -eq 1 ]; then
  exit 1
fi

echo "➜ 开始构建镜像并启动服务..."
$DC -f "$COMPOSE_FILE" up -d --build

echo
echo "✔ 服务部署完成："
$DC -f "$COMPOSE_FILE" ps
