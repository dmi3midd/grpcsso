package server

import (
	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/service"
)

type Server struct {
	grpcssov1.UnimplementedAuthServiceServer
	grpcssov1.UnimplementedResetPasswordServiceServer
	grpcssov1.UnimplementedPermissionServiceServer
	userService       service.UserService
	permissionService service.PermissionService
	resetService      service.ResetService
}

func NewServer(
	userService service.UserService,
	permissionService service.PermissionService,
	resetService service.ResetService,
) *Server {
	return &Server{
		userService:       userService,
		permissionService: permissionService,
		resetService:      resetService,
	}
}
