#!/bin/bash

protoc --go-grpc_out=. \
     --go_out=. \
     -I api/protobuf api/protobuf/*proto
