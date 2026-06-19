package server

import "github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"

type Server struct {
	grpcssov1.UnimplementedAuthServiceServer
}

func NewServer() *Server {
	return &Server{}
}
