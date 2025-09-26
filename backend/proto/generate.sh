#!/bin/bash

# 删除老文件
rm -rf ./pb/*.pb.go

# 生成新文件
protoc -I=./ \
  --go_out=./pb --go_opt=paths=source_relative \
  --go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=./pb --grpc-gateway_opt=paths=source_relative \
  *.proto
