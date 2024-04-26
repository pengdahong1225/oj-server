#!/bin/bash

# go
protoc --go_out=. --go-grpc_out=. *.proto

install *.go ../question-service/logic/proto
install *.go ../db-service/internal/proto

rm -rf *.go

# cpp
