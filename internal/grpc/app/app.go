package app

import (
	"fmt"
	"net"

	"github.com/dmi3midd/grpcsso/internal/grpc/apierrors"
	"github.com/dmi3midd/grpcsso/internal/grpc/interceptor"
	"github.com/dmi3midd/grpcsso/internal/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	gRPCServer *grpc.Server
}

func NewApp(srv *server.Server, opts ...grpc.ServerOption) (*App, error) {
	valInterceptor, err := interceptor.NewValidationInterceptor()
	if err != nil {
		return nil, fmt.Errorf("failed to create validation interceptor: %w", err)
	}

	errInterceptor := apierrors.UnaryErrorInterceptor()

	serverOpts := append([]grpc.ServerOption{
		grpc.ChainUnaryInterceptor(valInterceptor, errInterceptor),
	}, opts...)

	gRPCServer := grpc.NewServer(serverOpts...)
	reflection.Register(gRPCServer)
	srv.Register(gRPCServer)

	return &App{
		gRPCServer: gRPCServer,
	}, nil
}

func (a *App) Run(lis net.Listener) error {
	const op = "App.Run"

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	a.gRPCServer.GracefulStop()
}
