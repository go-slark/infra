package redis

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisClient *redis.Client

type RedisClientConfig struct {
	Address            string `json:"address"`
	Password           string `json:"password"`
	DB                 int    `json:"db"`
	DialTimeout        int    `json:"dial_timeout"`
	ReadTimeout        int    `json:"read_timeout"`
	WriteTimeout       int    `json:"write_timeout"`
	IdleTimeout        int    `json:"idle_timeout"`
	PoolTimeout        int    `json:"pool_timeout"`
	MaxConnAge         int    `json:"max_conn_age"`
	MaxRetry           int    `json:"max_retry"`
	PoolSize           int    `json:"pool_size"`
	MinIdleConns       int    `json:"min_idle_conns"`
	IdleCheckFrequency int    `json:"idle_check_frequency"`
	MaxRetryBackoff    int    `json:"max_retry_backoff"`
}

func InitRedisClient(c *RedisClientConfig) {
	client, err := createRedisClient(c)
	if err != nil {
		panic(errors.New(fmt.Sprintf("redis client %+v error %v", c, err)))
	}
	redisClient = client

}

func AppendRedisClients(config *RedisClientConfig) {
	if redisClient == nil {
		InitRedisClient(config)
	}

	client, err := createRedisClient(config)
	if err != nil {
		panic(errors.New(fmt.Sprintf("redis client %+v error %v", config, err)))
	}
	redisClient = client
}

func createRedisClient(c *RedisClientConfig) (*redis.Client, error) {
	options := &redis.Options{
		Network:  "tcp",
		Addr:     c.Address,
		Password: c.Password,
		DB:       c.DB,
	}

	if c.DialTimeout != 0 {
		options.DialTimeout = time.Duration(c.DialTimeout) * time.Second
	}
	if c.ReadTimeout != 0 {
		options.ReadTimeout = time.Duration(c.ReadTimeout) * time.Second
	}
	if c.WriteTimeout != 0 {
		options.WriteTimeout = time.Duration(c.WriteTimeout) * time.Second
	}
	if c.PoolTimeout != 0 {
		options.PoolTimeout = time.Duration(c.PoolTimeout) * time.Second
	}
	if c.MaxRetry != 0 {
		options.MaxRetries = c.MaxRetry
	}
	if c.PoolSize != 0 {
		options.PoolSize = c.PoolSize
	}
	if c.MinIdleConns != 0 {
		options.MinIdleConns = c.MinIdleConns
	}
	if c.MaxRetryBackoff != 0 {
		options.MaxRetryBackoff = time.Duration(c.MaxRetryBackoff) * time.Millisecond
	}
	client := redis.NewClient(options)
	_, err := client.Ping(context.TODO()).Result()
	return client, err
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func CloseRedisClients() error {
	if redisClient == nil {
		return nil
	}

	return redisClient.Close()
}
