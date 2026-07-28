package interceptor

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// NewValidationInterceptor returns a grpc.UnaryServerInterceptor that validates incoming requests
func NewValidationInterceptor() (grpc.UnaryServerInterceptor, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize protovalidate validator: %w", err)
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if msg, ok := req.(proto.Message); ok && msg != nil {
			if err := validator.Validate(msg); err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		return handler(ctx, req)
	}, nil
}
