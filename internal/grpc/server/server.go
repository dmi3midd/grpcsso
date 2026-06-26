package server

import "github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"

type Server struct {
	grpcssov1.UnimplementedAuthServiceServer
	grpcssov1.UnimplementedResetPasswordServiceServer
	grpcssov1.UnimplementedPermissionServiceServer
}

func NewServer() *Server {
	return &Server{}
}
