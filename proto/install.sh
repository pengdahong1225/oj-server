#!/bin/bash

protoc --go_out=. --go-grpc_out=. *.proto

install *.go ../user-service/proto
install *.go ../web-service/proto
install *.go ../cms-service/proto

rm -rf *.go