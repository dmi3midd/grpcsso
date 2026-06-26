package server

import (
	"context"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
)

func (s *Server) Registration(ctx context.Context, req *grpcssov1.RegistrationRequest) (*grpcssov1.RegistrationResponse, error) {
	return &grpcssov1.RegistrationResponse{}, nil
}

func (s *Server) Login(ctx context.Context, req *grpcssov1.LoginRequest) (*grpcssov1.LoginResponse, error) {
	return &grpcssov1.LoginResponse{}, nil
}

func (s *Server) Logout(ctx context.Context, req *grpcssov1.LogoutRequest) (*grpcssov1.LogoutResponse, error) {
	return &grpcssov1.LogoutResponse{}, nil
}

func (s *Server) Refresh(ctx context.Context, req *grpcssov1.RefreshRequest) (*grpcssov1.RefreshResponse, error) {
	return &grpcssov1.RefreshResponse{}, nil
}
