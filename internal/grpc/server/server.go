package server

import (
	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/service"
)

type Server struct {
	grpcssov1.UnimplementedAuthServiceServer
	grpcssov1.UnimplementedRBACServiceServer
	userService service.UserService
	rbacService service.RBACService
}

func NewServer() *Server {
	return &Server{}
}
