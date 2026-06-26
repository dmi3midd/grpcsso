package server

import (
	"context"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
)

func (s *Server) InitiateResetPassword(ctx context.Context, req *grpcssov1.InitiateResetPasswordRequest) (*grpcssov1.InitiateResetPasswordResponse, error) {
	return &grpcssov1.InitiateResetPasswordResponse{}, nil
}

func (s *Server) ConfirmResetPassword(ctx context.Context, req *grpcssov1.ConfirmResetPasswordRequest) (*grpcssov1.ConfirmResetPasswordResponse, error) {
	return &grpcssov1.ConfirmResetPasswordResponse{}, nil
}
