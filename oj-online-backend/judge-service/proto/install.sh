#!/bin/bash

# cpp
protoc --cpp_out=. --grpc_out=. *.proto
