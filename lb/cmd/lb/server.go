package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Run(ctx context.Context, cfg Config, out io.Writer) {
	configureLogger(out)

	server := newServer(cfg)

	go shutdownHandler(ctx, server)

	log.Printf("I'm Listening on port: %d\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Printf("Failed to listen and serve: %s\n", err)
	}
}

func configureLogger(out io.Writer) {
	log.SetPrefix("LB   ")
	log.SetOutput(out)
	log.SetFlags(log.Ltime | log.Lmsgprefix)
}

func newServer(cfg Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", requestHandler())
	webServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}
	return webServer
}

func shutdownHandler(
	ctx context.Context,
	server *http.Server,
) {
	<-ctx.Done()

	log.Println("Shutting down...")
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown cleanly: %s\n", err)
		os.Exit(1)
	}
}
