package redisps

import "github.com/gomodule/redigo/redis"

type Pool interface {
	Get() redis.Conn
}
