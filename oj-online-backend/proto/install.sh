#!/bin/bash

protoc --go_out=. --go-grpc_out=. *.proto

install *.go ../question-service/logic/proto
install *.go ../db-service/internal/proto
#install *.go ../judge-service/proto

rm -rf *.go