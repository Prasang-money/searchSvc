package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Prasang-money/searchSvc/route"
)

func main() {
	router := route.GetRoute()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Running server in a goroutine so that it doesn't block graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	// Setting up channel to listen for interrupt or terminate signal from OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// wait for interrupt signal
	<-ctx.Done()
	stop()

	log.Println("shutting down server gracefully")

	// The context is used to inform the server it has 10 seconds to finish
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server closed")
}
