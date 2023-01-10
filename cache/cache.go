package cache

import (
	"context"
	"time"
)

//go:generate mockgen --destination=./mocks/mock_cache.go --package=mocks github.com/PereRohit/util/cache Cacheable
type Cacheable interface {
	Health() error

	Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}
