// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package server

import (
	"github.com/simukti/grpc-password/proto"
	"github.com/simukti/grpc-password/storage"
	"google.golang.org/grpc"
)

type grpcServer struct {
	db storage.Storage
}

// New create new gRPC server
func New(db storage.Storage) *grpc.Server {
	rpc := &grpcServer{db: db}
	srv := grpc.NewServer()
	proto.RegisterPasswordServer(srv, rpc.PasswordServer())

	return srv
}
