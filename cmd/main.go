package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dmi3midd/grpcsso/internal/config"
	"github.com/dmi3midd/grpcsso/internal/grpc/app"
	"github.com/dmi3midd/grpcsso/internal/grpc/listener"
	"github.com/dmi3midd/grpcsso/internal/grpc/server"
	"github.com/dmi3midd/grpcsso/internal/postgres"
	"github.com/dmi3midd/grpcsso/internal/redis"
)

func main() {
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

	// Connect to Redis
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

	// Create gRPC server
	gRPCServer := server.NewServer()

	// Initialize gRPC app
	gRPCApp := app.NewApp(gRPCServer)

	// Run gRPC server in a goroutine
	go func() {
		log.Printf("starting gRPC server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := gRPCApp.Run(lis); err != nil {
			log.Fatalf("gRPC server failed to run: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Printf("received shutdown signal: %s", sign.String())

	gRPCApp.Stop()
	log.Println("application stopped gracefully")
}
