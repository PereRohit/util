package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	rdbCli redis.Cmdable
	cfg    *redis.Options
}

type RedisConfig struct {
	// host:port address.
	Address string `json:"address,omitempty"`
	// Use the specified Username to authenticate the current connection
	Username string `json:"username,omitempty"`
	// Optional password.
	Password string `json:"password,omitempty"`
	// Database to be selected after connecting to the server.
	DB int `json:"db,omitempty"`
	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int `json:"max_retries,omitempty"`
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration `json:"min_retry_backoff,omitempty"`
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration `json:"max_retry_backoff,omitempty"`
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration `json:"dial_timeout,omitempty"`
	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that fifo has higher overhead compared to lifo.
	PoolFIFO bool `json:"pool_fifo,omitempty"`
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int `json:"pool_size,omitempty"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	PoolTimeout time.Duration `json:"pool_timeout,omitempty"`
}

func NewRedis(cfg *RedisConfig) Cacheable {
	redisCfg := &redis.Options{
		Addr:            cfg.Address,
		Username:        cfg.Username,
		Password:        cfg.Password,
		DB:              cfg.DB,
		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
		DialTimeout:     cfg.DialTimeout,
		//PoolFIFO:           false,
		PoolSize:    cfg.PoolSize,
		PoolTimeout: cfg.PoolTimeout,
	}
	rdbCli := redis.NewClient(redisCfg)

	return &redisCache{
		cfg:    redisCfg,
		rdbCli: rdbCli,
	}
}

func (r redisCache) Health() error {
	return r.rdbCli.Ping(context.Background()).Err()
}

func (r redisCache) Set(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	return r.rdbCli.Set(ctx, key, value, expiry).Err()
}

func (r redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.rdbCli.Get(ctx, key).Bytes()
}

func (r redisCache) Delete(ctx context.Context, key string) error {
	return r.rdbCli.Del(ctx, key).Err()
}
