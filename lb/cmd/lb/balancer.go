package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

func requestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)

		address := selectService()
		resp, err := http.Get(address)
		if err != nil {
			msg := fmt.Sprintf("Failed to get response from backend: %s", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		defer closeResponse(resp.Body)

		text, err := io.ReadAll(resp.Body)
		if err != nil {
			msg := fmt.Sprintf("Could not read backend response: %s", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		log.Printf("Response from server: %s %s\n", resp.Proto, resp.Status)
		log.Printf("%s\n", string(text))

		if _, err := w.Write(text); err != nil {
			log.Printf("Failed to write response: %s\n", err)
		}
	})
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

func selectService() string {
	service := rand.Intn(3) + 1
	return fmt.Sprintf("http://be-%d", service)
}
