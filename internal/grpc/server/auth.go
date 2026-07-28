package server

import (
	"context"
	"errors"

	"github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func (s *Server) Registration(ctx context.Context, req *grpcssov1.RegistrationRequest) (*grpcssov1.RegistrationResponse, error) {
	if req.GetUsername() == "" || req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email and password are required")
	}

	id, err := s.userService.Registration(ctx, req.GetUsername(), req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExist) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.RegistrationResponse{
		UserId: id,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *grpcssov1.LoginRequest) (*grpcssov1.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	userAgent, ipAddress := extractMetadata(ctx)

	authDto, err := s.userService.Login(ctx, req.GetEmail(), req.GetPassword(), userAgent, ipAddress)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if errors.Is(err, service.ErrInvalidPassword) {
			return nil, status.Error(codes.Unauthenticated, "invalid password")
		}
		return nil, status.Error(codes.Internal, err.Error())
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
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	if err := s.userService.Logout(ctx, req.GetRefreshToken()); err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpcssov1.LogoutResponse{}, nil
}

func (s *Server) Refresh(ctx context.Context, req *grpcssov1.RefreshRequest) (*grpcssov1.RefreshResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	userAgent, ipAddress := extractMetadata(ctx)

	authDto, err := s.userService.Refresh(ctx, req.GetRefreshToken(), ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}
		if errors.Is(err, service.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
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
