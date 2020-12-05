package redisx

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

// CreatePool create redis pool
func CreatePool(address string, opts ...Option) *redis.Pool {
	var options = defaultPoolOptions

	for _, opt := range opts {
		opt(options)
	}

	return &redis.Pool{
		MaxIdle:     options.MaxIdle,
		MaxActive:   options.MaxActive,
		IdleTimeout: time.Duration(options.IdleTimeout) * time.Second,
		Wait:        options.Wait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address, poolDialOptions(options)...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: nil,
	}
}

func poolDialOptions(option *Options) []redis.DialOption {
	dialOptions := []redis.DialOption{
		redis.DialConnectTimeout(time.Duration(option.ConnectTimeout) * time.Second),
		redis.DialConnectTimeout(time.Duration(option.ReadTimeout) * time.Second),
		redis.DialWriteTimeout(time.Duration(option.WriteTimeout) * time.Second),
		redis.DialDatabase(option.DB),
	}

	if option.Password != "" {
		dialOptions = append(dialOptions, redis.DialPassword(option.Password))
	}
	return dialOptions
}
