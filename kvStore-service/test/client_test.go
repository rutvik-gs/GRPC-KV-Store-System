package test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "GRPC-KV-Store-System/schemas"
)

func TestKVStoreIntegration(t *testing.T) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	defer conn.Close()

	client := pb.NewKeyValueStoreClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("Set a key", func(t *testing.T) {
		resp, err := client.Set(ctx, &pb.SetRequest{
			Key:   "testkey",
			Value: "testvalue",
		})
		if err != nil {
			t.Fatalf("Set failed: %v", err)
		}
		if resp.Message == "" {
			t.Error("Expected success message")
		}
		t.Logf("Set response: %s", resp.Message)
	})

	t.Run("Retrieve a key", func(t *testing.T) {
		resp, err := client.Get(ctx, &pb.GetRequest{
			Key: "testkey",
		})
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if resp.Value != "testvalue" {
			t.Errorf("Expected 'testvalue', got '%s'", resp.Value)
		}
		t.Logf("Get response: %s", resp.Value)
	})

	t.Run("Retrieve a non existent key", func(t *testing.T) {
		_, err := client.Get(ctx, &pb.GetRequest{
			Key: "nonexistent",
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatal("Expected gRPC status error")
		}
		if st.Code() != codes.NotFound {
			t.Errorf("Expected NotFound, got %v", st.Code())
		}
		t.Logf("Expected error received: %s", st.Message())
	})

	t.Run("Try setting an empty key", func(t *testing.T) {
		_, err := client.Set(ctx, &pb.SetRequest{
			Key:   "",
			Value: "value",
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatal("Expected gRPC status error")
		}
		if st.Code() != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument, got %v", st.Code())
		}
		t.Logf("Expected error received: %s", st.Message())
	})

	t.Run("Delete a key", func(t *testing.T) {
		resp, err := client.Delete(ctx, &pb.DeleteRequest{
			Key: "testkey",
		})
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}
		if resp.Message == "" {
			t.Error("Expected success message")
		}
		t.Logf("Delete response: %s", resp.Message)
	})

	t.Run("Try retrieving a deleted key", func(t *testing.T) {
		_, err := client.Get(ctx, &pb.GetRequest{
			Key: "testkey",
		})
		if err == nil {
			t.Fatal("Expected error after delete, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatal("Expected gRPC status error")
		}
		if st.Code() != codes.NotFound {
			t.Errorf("Expected NotFound, got %v", st.Code())
		}
		t.Logf("Expected error received: %s", st.Message())
	})

	t.Run("Delete a non existent key", func(t *testing.T) {
		_, err := client.Delete(ctx, &pb.DeleteRequest{
			Key: "nonexistent",
		})
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatal("Expected gRPC status error")
		}
		if st.Code() != codes.NotFound {
			t.Errorf("Expected NotFound, got %v", st.Code())
		}
		t.Logf("Expected error received: %s", st.Message())
	})
}
