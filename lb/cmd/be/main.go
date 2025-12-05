package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	name := rand.Text()[:4]
	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt)
	ctx, cancel := signal.NotifyContext(ctx, os.Kill)
	defer cancel()

	if err := Run(
		ctx,
		name,
		os.Stdout,
	); err != nil {
		fmt.Printf("ERROR - exiting: %v\n", err)
		os.Exit(1)
	}
}

func Run(ctx context.Context, name string, out io.Writer) error {
	logger := slog.New(slog.NewTextHandler(out, nil)).With("service", name)
	slog.SetDefault(logger)

	port := 8080
	mux := http.NewServeMux()
	mux.Handle("/", requestHandler(name))
	webServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go shutdownHandler(ctx, webServer)

	slog.Info("I'm Listening", "port", port)
	if err := webServer.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		slog.Error("Failed to listen and serve", "error", err)
	}

	return nil
}

func shutdownHandler(
	ctx context.Context,
	server *http.Server,
) {
	<-ctx.Done()

	slog.Info("Shutting down...")
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Info("Failed to shutdown cleanly", "error", err)
		os.Exit(1)
	}
}

func requestHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "method", r.Method, "uri", r.RequestURI, "srvr", name)
		if _, err := fmt.Fprintf(w, "Hello from back end server: %q", name); err != nil {
			slog.Error("Failed to write response", "error", err)
		}
	})
}
