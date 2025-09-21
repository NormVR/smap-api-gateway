package main

import (
	"api-gateway/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	application := app.New()

	serverErr := make(chan error, 1)

	go application.Run(serverErr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
		log.Println("Received shutdown signal")
	case err := <-serverErr:
		log.Printf("Server error: %v\n", err)
	}

	log.Println("Shutting down services")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	application.Stop(ctx)
}
