#!/bin/sh

set -e

cd "$(dirname "$0")"

COMPOSE_FILE=docker-compose.env.yaml

# ---- 选择 docker compose 命令：v2 插件优先，兼容 v1 独立二进制 ----
if docker compose version >/dev/null 2>&1; then
  DC="docker compose"
elif command -v docker-compose >/dev/null 2>&1; then
  DC="docker-compose"
else
  echo "✘ 未检测到 docker compose，请先安装 docker 及 compose 插件" >&2
  exit 1
fi

echo "➜ 开始部署环境依赖..."
$DC -f "$COMPOSE_FILE" up -d

echo
echo "✔ 环境依赖部署完成："
$DC -f "$COMPOSE_FILE" ps

