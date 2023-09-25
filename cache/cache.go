package cache

import "context"

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, into interface{}) error
	Delete(ctx context.Context, key string) error
}
