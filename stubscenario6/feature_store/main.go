package main

import (
	"feature_store/srv"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", srv.Middleware(srv.Health))
	mux.HandleFunc("POST /execute", srv.Middleware(srv.Execute))

	port := 8000

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Starting '%s' on https://localhost:%d\n", srv.SrvName, port)
	log.Fatal(server.ListenAndServe())
}
