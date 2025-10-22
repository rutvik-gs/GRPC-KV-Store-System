package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"GRPC-KV-Store-System/kvStore-service/internal/server"
	"GRPC-KV-Store-System/kvStore-service/internal/store"
	pb "GRPC-KV-Store-System/schemas/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()

	log.Println("Starting gRPC server...")

	lis, ListenerErr := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if ListenerErr != nil {
		log.Fatalf("Failed to listen: %v", ListenerErr)
	}

	grpcServer := grpc.NewServer()

	kvStore := store.CreateStore()
	kvServer := server.StartServer(kvStore)

	pb.RegisterKeyValueStoreServer(grpcServer, kvServer)
	reflection.Register(grpcServer)

	log.Printf("gRPC server is now listening on port %d", *port)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	ServeErr := grpcServer.Serve(lis)

	if ServeErr != nil {
		log.Fatalf("Failed to serve: %v", ServeErr)
	}
}
