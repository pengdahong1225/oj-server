#!/bin/bash

protoc --go_out=. --go-grpc_out=. *.proto
protoc --cpp_out=. judge.proto

install *.go ../question-service/api/proto
install *.go ../db-service/internal/proto
rm -rf *.go