package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Run(ctx context.Context, cfg Config, out io.Writer) {
	configureLogger(cfg, out)

	server := newServer(cfg)

	go shutdownHandler(ctx, server)

	log.Printf("I'm Listening on port: %d\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Printf("Failed to listen and serve: %s\n", err)
	}
}

func configureLogger(cfg Config, out io.Writer) {
	log.SetPrefix(cfg.Name + " ")
	log.SetOutput(out)
	log.SetFlags(log.Ltime | log.Lmsgprefix)
}

func newServer(cfg Config) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", requestHandler(cfg.Name))
	webServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}
	return webServer
}

func requestHandler(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request from %s\n", r.Host)
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		for k, v := range r.Header {
			log.Printf("%s: %s", k, strings.Join(v, ","))
		}
		if _, err := fmt.Fprintf(w, "Hello from back end server: %q", name); err != nil {
			log.Printf("Failed to write response: %s\n", err)
		}
	})
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
