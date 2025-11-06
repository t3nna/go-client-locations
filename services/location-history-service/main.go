package main

import (
	"context"
	"go-clinet-locations/shared/db"
	"go-clinet-locations/shared/env"
	"go-clinet-locations/shared/messaging"
	grpcserver "google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var GrpcAddr = ":9092"

func main() {
	rabbitMqUri := env.GetString("RABBITMQ_URI",
		"amqp://guest:guest@rabbitmq:5672/")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Initialize MongoDB
	mongoClient, err := db.NewMongoClient(ctx, db.NewMongoDefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB, err: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	mongoDb := db.GetDatabase(mongoClient, db.NewMongoDefaultConfig())
	mongoDbRepo := NewMongoService(mongoDb)

	log.Printf(mongoDb.Name())

	// Rabbit mq
	conn, err := messaging.NewRabbitMQ(rabbitMqUri)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	//svc := NewService()
	//// starting the grpcServer
	grpcServer := grpcserver.NewServer()
	NewGrpcHandler(grpcServer, mongoDbRepo)

	log.Println("Starting gRPC server Location service on port ", lis.Addr().String())

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
