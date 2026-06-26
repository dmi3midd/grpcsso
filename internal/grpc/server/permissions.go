package server

import (
	"context"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
)

func (s *Server) HasPermissions(ctx context.Context, req *grpcssov1.HasPermissionsRequest) (*grpcssov1.HasPermissionsResponse, error) {
	return &grpcssov1.HasPermissionsResponse{}, nil
}

func (s *Server) GetPermissions(ctx context.Context, req *grpcssov1.GetPermissionsRequest) (*grpcssov1.GetPermissionsResponse, error) {
	return &grpcssov1.GetPermissionsResponse{}, nil
}

func (s *Server) AddPermissions(ctx context.Context, req *grpcssov1.AddPermissionsRequest) (*grpcssov1.AddPermissionsResponse, error) {
	return &grpcssov1.AddPermissionsResponse{}, nil
}

func (s *Server) RemovePermissions(ctx context.Context, req *grpcssov1.RemovePermissionsRequest) (*grpcssov1.RemovePermissionsResponse, error) {
	return &grpcssov1.RemovePermissionsResponse{}, nil
}
