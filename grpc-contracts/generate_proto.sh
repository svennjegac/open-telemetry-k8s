#!/bin/bash

# set script directory as current directory (for relative paths)
cd $(dirname $0)

# generate proto file
protoc --proto_path=. --go_out=plugins=grpc:.  ip_service.proto

# cleanup client protobuf generated file
if [ -d "../apps/user-http/src/internal/ip" ]; then
  rm -Rf ../apps/user-http/src/internal/ip;
fi

# create client protobuf generated file
mkdir ../apps/user-http/src/internal/ip
cp ip/ip_service.pb.go ../apps/user-http/src/internal/ip/ip_service.pb.go

# cleanup server protobuf generated file
if [ -d "../apps/ip-grpc/src/internal/ip" ]; then
  rm -Rf ../apps/ip-grpc/src/internal/ip;
fi

# create server protobuf generated file
mkdir ../apps/ip-grpc/src/internal/ip
cp ip/ip_service.pb.go ../apps/ip-grpc/src/internal/ip/ip_service.pb.go

# cleanup current directory unnecessary protobuf generated file
rm -Rf ip
