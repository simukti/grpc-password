#!/bin/sh
#protoc --go_out=plugins=grpc:. *.proto
protoc --gogoslick_out=plugins=grpc:. -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/gogo/protobuf/protobuf *.proto