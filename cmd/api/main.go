package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmi3midd/grpcsso/internal/config"
	"github.com/dmi3midd/grpcsso/internal/grpc/app"
	"github.com/dmi3midd/grpcsso/internal/grpc/listener"
	"github.com/dmi3midd/grpcsso/internal/grpc/server"
	"github.com/dmi3midd/grpcsso/internal/postgres"
	"github.com/dmi3midd/grpcsso/internal/redis"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/dmi3midd/grpcsso/internal/service"
	"github.com/dmi3midd/grpcsso/internal/workers"
)

func main() {
	// Root context with signal cancellation for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("starting application...")

	// Connect to Postgres
	postgresService, err := postgres.New(&cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer func() {
		log.Println("closing postgres connection")
		postgresService.Close()
	}()

	redisService, err := redis.New(&cfg.Redis)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer func() {
		log.Println("closing redis connection")
		redisService.Close()
	}()

	// Create listener
	listener := listener.NewListener(&cfg.Server)
	lis, err := listener.Listen()
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	// Initialize repositories and services
	db := postgresService.GetDB()
	txManager := repository.NewTxManager(db)
	redisClient := redisService.GetClient()

	userRepo := repository.NewUserRepo(db)
	tokenRepo := repository.NewTokenRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	permissionRepo := repository.NewPermissionRepo(db)
	roleCache := repository.NewRoleCache(redisClient)
	permissionCache := repository.NewPermissionCache(redisClient)

	cleanerInterval := 16 * time.Minute
	if cfg.JWT.TokenCleanerInterval > 0 {
		cleanerInterval = cfg.JWT.TokenCleanerInterval
	}
	tokenCleaner := workers.NewTokenCleaner(cleanerInterval, tokenRepo)
	go tokenCleaner.Start(ctx)

	tokenManager := service.NewTokenManager(txManager, tokenRepo, &cfg.Keys, &cfg.JWT)
	userService := service.NewUserService(txManager, userRepo, tokenManager)
	rbacService := service.NewRBACService(txManager, userRepo, roleRepo, permissionRepo, roleCache, permissionCache)

	// Create gRPC server
	gRPCServer := server.NewServer(userService, rbacService)

	// Initialize gRPC app
	gRPCApp, err := app.NewApp(gRPCServer)
	if err != nil {
		log.Fatalf("failed to create gRPC app: %v", err)
	}

	// Run gRPC server in a goroutine
	go func() {
		log.Printf("starting gRPC server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := gRPCApp.Run(lis); err != nil {
			log.Fatalf("gRPC server failed to run: %v", err)
		}
	}()

	// Graceful shutdown
	<-ctx.Done()
	log.Println("received shutdown signal, stopping application...")

	gRPCApp.Stop()
	log.Println("application stopped gracefully")
}
