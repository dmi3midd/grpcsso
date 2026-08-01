package apierrors

import (
	"context"
	"log/slog"

	"errors"

	"github.com/dmi3midd/grpcsso/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryErrorInterceptor returns a gRPC UnaryServerInterceptor that maps domain errors to gRPC status codes.
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, HandleError(ctx, req, info, err)
		}
		return resp, nil
	}
}

func HandleError(ctx context.Context, req any, info *grpc.UnaryServerInfo, err error) error {

	if err == nil {
		return nil
	}

	// If error is already a gRPC status error, return it as is
	if st, ok := status.FromError(err); ok && st.Code() != codes.Unknown {
		return err
	}

	slog.Error(
		"error",
		slog.String("method", info.FullMethod),
		slog.String("error", err.Error()),
	)

	switch {

	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")

	case errors.Is(err, service.ErrUserAlreadyExist):
		return status.Error(codes.AlreadyExists, "user already exists")

	case errors.Is(err, service.ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, "invalid password")

	case errors.Is(err, service.ErrInvalidRefreshToken):
		return status.Error(codes.InvalidArgument, "invalid refresh token")

	case errors.Is(err, service.ErrInvalidAccessToken):
		return status.Error(codes.InvalidArgument, "invalid access token")

	case errors.Is(err, service.ErrSubjectAndIdNotFound), errors.Is(err, service.ErrUnexpectedSigningMethod):
		return status.Error(codes.Internal, "internal error")

	case errors.Is(err, service.ErrTokenNotFound):
		return status.Error(codes.Unauthenticated, "refresh token expired or invalid")

	case errors.Is(err, service.ErrRoleNotFound):
		return status.Error(codes.NotFound, "role not found")

	case errors.Is(err, service.ErrPermissionNotFound):
		return status.Error(codes.NotFound, "permission not found")

	default:
		return status.Error(codes.Internal, "internal error")
	}
}
