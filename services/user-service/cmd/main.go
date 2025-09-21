package main

import (
	"context"
	"go-clinet-locations/services/user-service/internal/infrastructure/grpc"
	"go-clinet-locations/services/user-service/internal/infrastructure/repository"
	"go-clinet-locations/services/user-service/internal/service"
	"go-clinet-locations/shared/db"
	grpcserver "google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var GrpcAddr = ":9093"

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize MongoDB
	mongoClient, err := db.NewMongoClient(ctx, db.NewMongoDefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB, err: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	mongoDb := db.GetDatabase(mongoClient, db.NewMongoDefaultConfig())
	mongoDbRepo := repository.NewMongoRepository(mongoDb)

	log.Printf(mongoDb.Name())

	//inmemRepo := repository.NewInmemRepository()
	//svc := service.NewService(inmemRepo)
	svc := service.NewService(mongoDbRepo)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, svc)

	log.Println("Starting gRPC server Trip service on port ", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()

	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()

}
