package server

import (
	"context"
	"errors"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetPermission(ctx context.Context, req *grpcssov1.GetPermissionRequest) (*grpcssov1.GetPermissionResponse, error) {
	permission, err := s.rbacService.GetPermissionById(ctx, req.GetPermissionId())
	if err != nil {
		if errors.Is(err, service.ErrPermissionNotFound) {
			return nil, status.Error(codes.NotFound, "permission not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.GetPermissionResponse{
		PermissionId: permission.Id,
		Name:         permission.Name,
	}, nil
}

func (s *Server) CreatePermission(ctx context.Context, req *grpcssov1.CreatePermissionRequest) (*grpcssov1.CreatePermissionResponse, error) {
	id, err := s.rbacService.CreatePermission(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.CreatePermissionResponse{
		PermissionId: id,
	}, nil
}

func (s *Server) DeletePermission(ctx context.Context, req *grpcssov1.DeletePermissionRequest) (*grpcssov1.DeletePermissionResponse, error) {
	id, err := s.rbacService.DeletePermission(ctx, req.GetPermissionId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.DeletePermissionResponse{
		PermissionId: id,
	}, nil
}

func (s *Server) AssignPermission(ctx context.Context, req *grpcssov1.AssignPermissionRequest) (*grpcssov1.AssignPermissionResponse, error) {
	permId, roleId, err := s.rbacService.AssignPermissionToRole(ctx, req.GetPermissionId(), req.GetRoleId())
	if err != nil {
		if errors.Is(err, service.ErrPermissionNotFound) {
			return nil, status.Error(codes.NotFound, "permission not found")
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.AssignPermissionResponse{
		PermissionId: permId,
		RoleId:       roleId,
	}, nil
}

func (s *Server) RevokePermission(ctx context.Context, req *grpcssov1.RevokePermissionRequest) (*grpcssov1.RevokePermissionResponse, error) {
	permId, roleId, err := s.rbacService.RevokePermissionFromRole(ctx, req.GetPermissionId(), req.GetRoleId())
	if err != nil {
		if errors.Is(err, service.ErrPermissionNotFound) {
			return nil, status.Error(codes.NotFound, "permission not found")
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.RevokePermissionResponse{
		PermissionId: permId,
		RoleId:       roleId,
	}, nil
}

func (s *Server) GetRolePermissions(ctx context.Context, req *grpcssov1.GetRolePermissionsRequest) (*grpcssov1.GetRolePermissionsResponse, error) {
	permissions, err := s.rbacService.GetRolePermissions(ctx, req.GetRoleId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoPerms := make([]*grpcssov1.Permission, 0, len(permissions))
	for _, p := range permissions {
		protoPerms = append(protoPerms, &grpcssov1.Permission{
			PermissionId: p.Id,
			Name:         p.Name,
		})
	}

	return &grpcssov1.GetRolePermissionsResponse{
		RoleId:      req.GetRoleId(),
		Permissions: protoPerms,
	}, nil
}
