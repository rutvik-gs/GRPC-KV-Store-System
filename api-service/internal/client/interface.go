package client

// ClientInterface defines the contract for KV store operations
type ClientInterface interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Close() error
}
