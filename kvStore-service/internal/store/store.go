package store

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("key cannot be empty")
)

type Store interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}
