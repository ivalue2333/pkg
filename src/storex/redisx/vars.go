package redisx

type (
	RedisType int
)


const (
	RedisTypeNode RedisType = 1
	RedisTypeCluster RedisType = 2
)

var (
	defaultPoolOptions = &Options{
		ConnectTimeout: 2,
		ReadTimeout:    2,
		WriteTimeout:   2,
		MaxActive:      200,
		MaxIdle:        200,
		IdleTimeout:    2,
		Wait:           false,
		DB:             0,
		Password:       "",
	}
)
