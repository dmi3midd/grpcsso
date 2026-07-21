package server

import (
	"context"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
)

func (s *Server) GetRole(ctx context.Context, req *grpcssov1.GetRoleRequest) (*grpcssov1.GetRoleResponse, error) {
	panic("implemnt getRole")
}

func (s *Server) CreateRole(ctx context.Context, req *grpcssov1.CreateRoleRequest) (*grpcssov1.CreateRoleResponse, error) {
	panic("implement createRole")
}

func (s *Server) DeleteRole(ctx context.Context, req *grpcssov1.DeleteRoleRequest) (*grpcssov1.DeleteRoleResponse, error) {
	panic("implement deleteRole")
}

func (s *Server) AssignRole(ctx context.Context, req *grpcssov1.AssignRoleRequest) (*grpcssov1.AssignRoleResponse, error) {
	panic("implement assignRole")
}
func (s *Server) RevokeRole(ctx context.Context, req *grpcssov1.RevokeRoleRequest) (*grpcssov1.RevokeRoleResponse, error) {
	panic("implement revokeRole")
}

func (s *Server) GetUserRoles(ctx context.Context, req *grpcssov1.GetUserRolesRequest) (*grpcssov1.GetUserRolesResponse, error) {
	panic("implement getUserRoles")
}
