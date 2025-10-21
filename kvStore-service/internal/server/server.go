package server

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GRPC-KV-Store-System/kvStore-service/internal/store"
	pb "GRPC-KV-Store-System/schemas"
)

type Server struct {
	pb.UnimplementedKeyValueStoreServer
	store store.Store
}

func StartServer(i store.Store) *Server {
	return &Server{
		store: i,
	}
}

func (i *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	log.Printf("Processing Request: Set key=%s, value=%s", req.Key, req.Value)

	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	err := i.store.Set(req.Key, req.Value)

	if err != nil {
		if errors.Is(err, store.ErrEmptyKey) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to store value: %v", err)
	}

	log.Printf("Successfully stored key=%s, value=%s", req.Key, req.Value)

	return &pb.SetResponse{
		Message: "Value Stored Successfully",
	}, nil
}

func (i *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Processing Request: Get key=%s", req.Key)

	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	value, err := i.store.Get(req.Key)

	if err != nil {
		if errors.Is(err, store.ErrEmptyKey) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, status.Error(codes.NotFound, "key not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to retrieve value: %v", err)
	}

	log.Printf("Successfully retrieved key=%s", req.Key)

	return &pb.GetResponse{
		Value: value,
	}, nil
}

func (i *Server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	log.Printf("Processing Request: Delete key=%s", req.Key)

	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	err := i.store.Delete(req.Key)
	if err != nil {
		if errors.Is(err, store.ErrEmptyKey) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if errors.Is(err, store.ErrKeyNotFound) {
			return nil, status.Error(codes.NotFound, "key not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to delete key: %v", err)
	}

	log.Printf("Successfully deleted key=%s", req.Key)

	return &pb.DeleteResponse{
		Message: "key deleted successfully",
	}, nil
}
