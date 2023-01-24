package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	redisClients = make(map[string]*redis.Client)
	redisOnce    sync.Once
)

type RedisClientConfig struct {
	Alias              string `json:"alias"`
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

func InitRedisClients(configs []*RedisClientConfig) error {
	redisOnce.Do(func() {
		for _, c := range configs {
			if _, ok := redisClients[c.Alias]; ok {
				panic(errors.New("duplicate redis client: " + c.Alias))
			}
			client, err := createRedisClient(c)
			if err != nil {
				panic(errors.New(fmt.Sprintf("redis client %+v error %v", c, err)))
			}
			redisClients[c.Alias] = client
		}
	})

	return nil
}

func AppendRedisClients(configs []*RedisClientConfig) {
	if len(redisClients) == 0 {
		_ = InitRedisClients(configs)
	}

	for _, c := range configs {
		if _, ok := redisClients[c.Alias]; ok {
			continue
		}
		client, err := createRedisClient(c)
		if err != nil {
			panic(errors.New(fmt.Sprintf("redis client %+v error %v", c, err)))
		}
		redisClients[c.Alias] = client
	}
}

func createRedisClient(c *RedisClientConfig) (*redis.Client, error) {
	options := &redis.Options{
		Network:     "tcp",
		Addr:        c.Address,
		Password:    c.Password,
		DB:          c.DB,
		IdleTimeout: time.Duration(c.IdleTimeout) * time.Second,
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
	if c.IdleTimeout != 0 {
		options.IdleTimeout = time.Duration(c.IdleTimeout) * time.Minute
	}
	if c.PoolTimeout != 0 {
		options.PoolTimeout = time.Duration(c.PoolTimeout) * time.Second
	}
	if c.MaxConnAge != 0 {
		options.MaxConnAge = time.Duration(c.MaxConnAge) * time.Second
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
	if c.IdleCheckFrequency != 0 {
		options.IdleCheckFrequency = time.Duration(c.IdleCheckFrequency) * time.Second
	}
	if c.MaxRetryBackoff != 0 {
		options.MaxRetryBackoff = time.Duration(c.MaxRetryBackoff) * time.Millisecond
	}
	redisClient := redis.NewClient(options)
	_, err := redisClient.Ping().Result()
	return redisClient, err
}

func GetRedisClient(alias string) *redis.Client {
	return redisClients[alias]
}

func CloseRedisClients() {
	for _, client := range redisClients {
		if client == nil {
			continue
		}

		_ = client.Close()
	}
}
