package app

import (
	"fmt"
	"net"

	grpcssov1 "github.com/dmi3midd/grpcsso-protos/gen/go/grpcssov1"
	"github.com/dmi3midd/grpcsso/internal/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	gRPCServer *grpc.Server
}

func NewApp(srv *server.Server) *App {
	gRPCServer := grpc.NewServer()

	grpcssov1.RegisterAuthServiceServer(gRPCServer, srv)
	grpcssov1.RegisterPermissionServiceServer(gRPCServer, srv)
	grpcssov1.RegisterResetPasswordServiceServer(gRPCServer, srv)
	// TODO: Enable reflection in production
	reflection.Register(gRPCServer)

	return &App{
		gRPCServer: gRPCServer,
	}
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
