#!/bin/bash

protoc --go_out=. --go-grpc_out=. *.proto

install *.go ../question-service/internal/proto
install *.go ../db-service/internal/proto
install *.go ../judge-service/internal/proto
rm -rf *.go