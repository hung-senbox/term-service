package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
}
