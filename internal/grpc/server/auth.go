package server

import (
	"context"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/grpc/apierrors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func (s *Server) Registration(ctx context.Context, req *grpcssov1.RegistrationRequest) (*grpcssov1.RegistrationResponse, error) {
	id, err := s.userService.Registration(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, apierrors.HandleError(ctx, req, nil, err)
	}

	return &grpcssov1.RegistrationResponse{
		UserId: id,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *grpcssov1.LoginRequest) (*grpcssov1.LoginResponse, error) {
	userAgent, ipAddress := extractMetadata(ctx)

	authDto, err := s.userService.Login(ctx, req.GetEmail(), req.GetPassword(), userAgent, ipAddress)
	if err != nil {
		return nil, apierrors.HandleError(ctx, req, nil, err)
	}

	return &grpcssov1.LoginResponse{
		User: &grpcssov1.User{
			UserId:   authDto.User.Id,
			Username: authDto.User.Username,
			Email:    authDto.User.Email,
		},
		AccessToken:  authDto.AccessToken,
		RefreshToken: authDto.RefreshToken,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *grpcssov1.LogoutRequest) (*grpcssov1.LogoutResponse, error) {
	if err := s.userService.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, apierrors.HandleError(ctx, req, nil, err)
	}

	return &grpcssov1.LogoutResponse{}, nil
}

func (s *Server) Refresh(ctx context.Context, req *grpcssov1.RefreshRequest) (*grpcssov1.RefreshResponse, error) {
	userAgent, ipAddress := extractMetadata(ctx)

	authDto, err := s.userService.Refresh(ctx, req.GetRefreshToken(), ipAddress, userAgent)
	if err != nil {
		return nil, apierrors.HandleError(ctx, req, nil, err)
	}

	return &grpcssov1.RefreshResponse{
		User: &grpcssov1.User{
			UserId:   authDto.User.Id,
			Username: authDto.User.Username,
			Email:    authDto.User.Email,
		},
		AccessToken:  authDto.AccessToken,
		RefreshToken: authDto.RefreshToken,
	}, nil
}

func extractMetadata(ctx context.Context) (userAgent, ipAddress string) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("user-agent"); len(ua) > 0 {
			userAgent = ua[0]
		}
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			ipAddress = xff[0]
		}
	}
	if ipAddress == "" {
		if pr, ok := peer.FromContext(ctx); ok && pr.Addr != nil {
			ipAddress = pr.Addr.String()
		}
	}
	return userAgent, ipAddress
}
