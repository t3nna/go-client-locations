package main

import (
	"context"
	"fmt"
	"go-clinet-locations/shared/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8004")
)

func main() {
	fmt.Printf("Starting API Gateway ")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/location", enableCORS(HandleUpdateUserLocation))

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server is listeneing on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the server: %v", err)
	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the server gracefully %v", err)
			server.Close()
		}
	}

}
