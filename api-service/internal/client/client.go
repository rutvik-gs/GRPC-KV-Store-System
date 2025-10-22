package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "GRPC-KV-Store-System/schemas/grpc"
)

type KVStoreClient struct {
	client pb.KeyValueStoreClient
	conn   *grpc.ClientConn
}

func NewKVStoreClient(grpcServerAddr string) (*KVStoreClient, error) {
	conn, err := grpc.NewClient(grpcServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	client := pb.NewKeyValueStoreClient(conn)
	log.Printf("Connected to gRPC server at %s", grpcServerAddr)

	return &KVStoreClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *KVStoreClient) Close() error {
	return c.conn.Close()
}

func (c *KVStoreClient) Set(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.client.Set(ctx, &pb.SetRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (c *KVStoreClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Get(ctx, &pb.GetRequest{
		Key: key,
	})
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

func (c *KVStoreClient) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.client.Delete(ctx, &pb.DeleteRequest{
		Key: key,
	})
	return err
}
