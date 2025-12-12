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
	b := NewBalancer(cfg.Urls)
	mux := http.NewServeMux()
	mux.Handle("/", requestHandler(b))
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

func requestHandler(b *Balancer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		address, err := b.selectService()
		if err != nil {
			msg := fmt.Sprintf("No backend services available: %s", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
		}

		text, err := forwardRequest(address)
		if err != nil {
			msg := fmt.Sprintf("Failed to get response from backend: %s", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(text); err != nil {
			log.Printf("Failed to write response: %s\n", err)
		}
	})
}

func forwardRequest(address string) ([]byte, error) {
	resp, err := http.Get(address)
	if err != nil {
		return nil, fmt.Errorf("getting response from backend: %w", err)
	}
	defer closeResponse(resp.Body)

	text, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from backend: %w", err)
	}

	log.Printf("Response from server: %s %s\n", resp.Proto, resp.Status)
	log.Printf("%s\n", string(text))
	return text, nil
}

func logRequest(r *http.Request) {
	log.Printf("Received request from %s\n", r.Host)
	log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
	for k, v := range r.Header {
		log.Printf("%s: %s", k, strings.Join(v, ","))
	}
	log.Println()
}

func closeResponse(r io.Closer) {
	if err := r.Close(); err != nil {
		log.Printf("Failed to close request: %s\n", err)
	}
}
