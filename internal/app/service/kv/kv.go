package kv

import (
	"context"
)

type KVStorage interface {
	Set(ctx context.Context, key, value string) (bool, error)
	Get(ctx context.Context, key string) (string, error)
}
