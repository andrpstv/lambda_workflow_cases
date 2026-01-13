package main

import (
	"model_4/srv"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", srv.Middleware(srv.Health))
	mux.HandleFunc("POST /execute", srv.Middleware(srv.Execute))

	port := 8060

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Starting '%s' on https://localhost:%d\n", srv.SrvName, port)
	log.Fatal(server.ListenAndServe())
}
