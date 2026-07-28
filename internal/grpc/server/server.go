package server

import (
	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/service"
	"google.golang.org/grpc"
)

type Server struct {
	grpcssov1.UnimplementedAuthServiceServer
	grpcssov1.UnimplementedRBACServiceServer
	userService service.UserService
	rbacService service.RBACService
}

func NewServer(userService service.UserService, rbacService service.RBACService) *Server {
	return &Server{
		userService: userService,
		rbacService: rbacService,
	}
}

func (s *Server) Register(gRPCServer *grpc.Server) {
	grpcssov1.RegisterAuthServiceServer(gRPCServer, s)
	grpcssov1.RegisterRBACServiceServer(gRPCServer, s)
}
