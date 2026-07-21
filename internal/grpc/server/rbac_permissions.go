package server

import (
	"context"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
)

func (s *Server) GetPermission(ctx context.Context, req *grpcssov1.GetPermissionRequest) (*grpcssov1.GetPermissionResponse, error) {
	panic("implement getPermission")
}

func (s *Server) CreatePermission(ctx context.Context, req *grpcssov1.CreatePermissionRequest) (*grpcssov1.CreatePermissionResponse, error) {
	panic("implement createPermission")
}

func (s *Server) DeletePermission(ctx context.Context, req *grpcssov1.DeletePermissionRequest) (*grpcssov1.DeletePermissionResponse, error) {
	panic("implement deletePermission")
}

func (s *Server) AssignPermission(ctx context.Context, req *grpcssov1.AssignPermissionRequest) (*grpcssov1.AssignPermissionResponse, error) {
	panic("implement assignPermission")
}
func (s *Server) RevokePermission(ctx context.Context, req *grpcssov1.RevokePermissionRequest) (*grpcssov1.RevokePermissionResponse, error) {
	panic("implement revokePermission")
}

func (s *Server) GetRolePermissions(ctx context.Context, req *grpcssov1.GetRolePermissionsRequest) (*grpcssov1.GetRolePermissionsResponse, error) {
	panic("implement getRolePermissions")
}
