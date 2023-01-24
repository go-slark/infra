package redis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var (
	redisPools = make(map[string]*redis.Pool)
	once       sync.Once
)

type RedisPoolConfig struct {
	Alias          string `json:"alias"`
	Address        string `json:"address"`
	Password       string `json:"password"`
	DB             *int   `json:"db"`
	ConnectTimeout int    `json:"connect_timeout"`
	ReadTimeout    int    `json:"read_timeout"`
	WriteTimeout   int    `json:"write_timeout"`
	Wait           bool   `json:"wait"`
	MaxIdle        int    `json:"max_idle"`
	IdleTimeout    int    `json:"idle_timeout"`
}

func InitRedisPool(configs []*RedisPoolConfig) error {
	once.Do(func() {
		for _, c := range configs {
			if _, ok := redisPools[c.Alias]; ok {
				panic(errors.New("duplicate redis pool: " + c.Alias))
			}
			p, err := createRedisPool(c)
			if err != nil {
				panic(errors.New(fmt.Sprintf("redis pool %+v error %v", c, err)))
			}
			redisPools[c.Alias] = p
		}
	})

	return nil
}

func createRedisPool(c *RedisPoolConfig) (*redis.Pool, error) {
	p := &redis.Pool{
		MaxIdle:     c.MaxIdle,
		IdleTimeout: time.Duration(c.IdleTimeout) * time.Second,
		Wait:        c.Wait,
		Dial: func() (conn redis.Conn, err error) {
			var options []redis.DialOption

			if c.ConnectTimeout != 0 {
				options = append(options, redis.DialConnectTimeout(time.Duration(c.ConnectTimeout)*time.Second))
			}
			if c.ReadTimeout != 0 {
				options = append(options, redis.DialReadTimeout(time.Duration(c.ReadTimeout)*time.Second))
			}
			if c.WriteTimeout != 0 {
				options = append(options, redis.DialWriteTimeout(time.Duration(c.WriteTimeout)*time.Second))
			}
			conn, err = redis.Dial(
				"tcp",
				c.Address,
				options...,
			)

			if err != nil {
				return nil, err
			}

			if c.Password != "" {
				if _, err := conn.Do("AUTH", c.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}

			if c.DB != nil {
				if _, err := conn.Do("SELECT", *c.DB); err != nil {
					conn.Close()
					return nil, err
				}
			}

			return conn, nil
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	conn := p.Get()
	_, err := conn.Do("PING")
	return p, err
}

func GetRedisPool(alias string) *redis.Pool {
	return redisPools[alias]
}

func WithConn(poolAlias string, fn func(conn redis.Conn) error) error {
	conn := GetRedisPool(poolAlias).Get()
	defer conn.Close()
	return fn(conn)
}

func Close() {
	for _, p := range redisPools {
		if p == nil {
			continue
		}

		_ = p.Close()
	}
}
