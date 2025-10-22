package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"GRPC-KV-Store-System/api-service/internal/client"
	"GRPC-KV-Store-System/api-service/internal/handler"
	"GRPC-KV-Store-System/api-service/internal/middleware"
)

var (
	port           = flag.String("port", "8080", "HTTP server port")
	grpcServerAddr = flag.String("grpc-addr", "localhost:50051", "gRPC server address")
	specPath       = flag.String("spec", "../schemas/rest/openapi.yaml", "OpenAPI spec path")
)

func main() {
	flag.Parse()

	if addr := os.Getenv("GRPC_SERVER_ENDPOINT"); addr != "" {
		*grpcServerAddr = addr
	}

	log.Println("Starting REST API server...")

	grpcClient, err := client.NewKVStoreClient(*grpcServerAddr)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	validator, err := middleware.NewValidationMiddleware(*specPath)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	h := handler.NewHandler(grpcClient)

	router := mux.NewRouter()

	router.HandleFunc("/health", h.HealthHandler).Methods("GET")
	router.HandleFunc("/kv", h.SetHandler).Methods("POST")
	router.HandleFunc("/kv/{key}", h.GetHandler).Methods("GET")
	router.HandleFunc("/kv/{key}", h.DeleteHandler).Methods("DELETE")

	router.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *specPath)
	}).Methods("GET")

	validatedRouter := validator.Validate(router)

	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      validatedRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("REST API server listening on port %s", *port)
		log.Printf("OpenAPI spec available at http://localhost:%s/openapi.yaml", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down REST API server...")
}
