#!/bin/bash

config_dir=".."
source_config="config"
compose_file="docker-compose.service.yml"

# 遍历当前目录下以-service结尾的服务目录
for service_dir in */
do
  # 排除当前目录和隐藏目录
  if [[ "$service_dir" != "." && "$service_dir" != ".." ]]; then
    # 截取服务目录名，去除尾随的斜杠
    service_name="${service_dir%/}"
    
    # 检查是否以-service结尾，以确定是否为服务目录
    if [[ $service_name == *"-service" ]]; then
      echo "Copying config to $service_name..."
      
      # 创建目标config文件夹（如果不存在）
      mkdir -p "$service_name/$source_config"
      
      # 将../config下的所有内容复制到每个服务目录的config文件夹中
      cp -r "$config_dir/$source_config/"* "$service_name/$source_config/"
      
      echo "Copied config to $service_name/$source_config/"
    fi
  fi
done

# 起服务
echo "Starting services with docker-compose..."
docker compose -f $compose_file up -d
