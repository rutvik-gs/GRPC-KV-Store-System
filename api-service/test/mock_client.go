package test

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GRPC-KV-Store-System/api-service/internal/client"
)

// MockClient is a test double for the gRPC client
type MockClient struct {
	store map[string]string
}

func NewMockClient() *MockClient {
	return &MockClient{
		store: make(map[string]string),
	}
}

func (m *MockClient) Set(key, value string) error {
	if key == "" {
		return status.Error(codes.InvalidArgument, "key cannot be empty")
	}
	m.store[key] = value
	return nil
}

func (m *MockClient) Get(key string) (string, error) {
	if key == "" {
		return "", status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	value, exists := m.store[key]
	if !exists {
		return "", status.Error(codes.NotFound, "key not found")
	}

	return value, nil
}

func (m *MockClient) Delete(key string) error {
	if key == "" {
		return status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	if _, exists := m.store[key]; !exists {
		return status.Error(codes.NotFound, "key not found")
	}

	delete(m.store, key)
	return nil
}

func (m *MockClient) Close() error {
	return nil
}

// Ensure MockClient implements ClientInterface
var _ client.ClientInterface = (*MockClient)(nil)
