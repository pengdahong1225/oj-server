#!/bin/bash

protoc --go_out=. --go-grpc_out=. *.proto

install *.go ../question-service/proto
install *.go ../db-service/proto
#install *.go ../judge-service/proto

rm -rf *.go