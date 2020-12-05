package redisx

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"time"
)

func NewCluster(startupNodes []string, opts ...Option) (*redisc.Cluster, error) {
	if len(startupNodes) == 0 {
		return nil, errors.New("no redis cluster startup nodes")
	}

	var options = defaultPoolOptions
	for _, opt := range opts {
		opt(options)
	}

	cluster := &redisc.Cluster{
		StartupNodes: startupNodes,
		DialOptions:  clusterDialOptions(options),
		CreatePool:   clusterCreatePool(options),
	}
	return cluster, cluster.Refresh()
}

func clusterCreatePool(options *Options) func(address string, opts ...redis.DialOption) (*redis.Pool, error) {
	return func(address string, opts ...redis.DialOption) (*redis.Pool, error) {
		opts = append(opts, clusterDialOptions(options)...)
		return &redis.Pool{
			MaxIdle:     options.MaxIdle,
			MaxActive:   options.MaxActive,
			IdleTimeout: time.Duration(options.IdleTimeout) * time.Second,
			Wait:        options.Wait,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", address, opts...)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: nil,
		}, nil
	}
}

func clusterDialOptions(option *Options) []redis.DialOption {
	dialOptions := []redis.DialOption{
		redis.DialConnectTimeout(time.Duration(option.ConnectTimeout) * time.Second),
		redis.DialConnectTimeout(time.Duration(option.ReadTimeout) * time.Second),
		redis.DialWriteTimeout(time.Duration(option.WriteTimeout) * time.Second),
	}
	if option.Password != "" {
		dialOptions = append(dialOptions, redis.DialPassword(option.Password))
	}
	return dialOptions
}
