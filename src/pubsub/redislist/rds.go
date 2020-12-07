package redislist

import "github.com/gomodule/redigo/redis"

type Pool interface {
	Get() redis.Conn
}
