package server

import (
	"context"
	"errors"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetRole(ctx context.Context, req *grpcssov1.GetRoleRequest) (*grpcssov1.GetRoleResponse, error) {
	role, err := s.rbacService.GetRoleById(ctx, req.GetRoleId())
	if err != nil {
		if errors.Is(err, service.ErrRoleNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.GetRoleResponse{
		RoleId: role.Id,
		Name:   role.Name,
	}, nil
}

func (s *Server) CreateRole(ctx context.Context, req *grpcssov1.CreateRoleRequest) (*grpcssov1.CreateRoleResponse, error) {
	id, err := s.rbacService.CreateRole(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.CreateRoleResponse{
		RoleId: id,
	}, nil
}

func (s *Server) DeleteRole(ctx context.Context, req *grpcssov1.DeleteRoleRequest) (*grpcssov1.DeleteRoleResponse, error) {
	id, err := s.rbacService.DeleteRole(ctx, req.GetRoleId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.DeleteRoleResponse{
		RoleId: id,
	}, nil
}

func (s *Server) AssignRole(ctx context.Context, req *grpcssov1.AssignRoleRequest) (*grpcssov1.AssignRoleResponse, error) {
	roleId, userId, err := s.rbacService.AssignRoleToUser(ctx, req.GetRoleId(), req.GetUserId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.AssignRoleResponse{
		RoleId: roleId,
		UserId: userId,
	}, nil
}

func (s *Server) RevokeRole(ctx context.Context, req *grpcssov1.RevokeRoleRequest) (*grpcssov1.RevokeRoleResponse, error) {
	roleId, userId, err := s.rbacService.RevokeRoleFromUser(ctx, req.GetRoleId(), req.GetUserId())
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, service.ErrRoleNotFound) {
			return nil, status.Error(codes.NotFound, "role not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.RevokeRoleResponse{
		RoleId: roleId,
		UserId: userId,
	}, nil
}

func (s *Server) GetUserRoles(ctx context.Context, req *grpcssov1.GetUserRolesRequest) (*grpcssov1.GetUserRolesResponse, error) {
	roles, err := s.rbacService.GetUserRoles(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoRoles := make([]*grpcssov1.Role, 0, len(roles))
	for _, r := range roles {
		protoRoles = append(protoRoles, &grpcssov1.Role{
			RoleId: r.Id,
			Name:   r.Name,
		})
	}

	return &grpcssov1.GetUserRolesResponse{
		UserId: req.GetUserId(),
		Roles:  protoRoles,
	}, nil
}
