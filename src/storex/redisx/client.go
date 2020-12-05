package redisx

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type (
	Client interface {
		GetConn() (conn Conn)
		DoWithContext(ctx context.Context, cmd string, args ...interface{}) (reply interface{}, err error)
		Close()
	}

	conns interface {
		Get() redis.Conn
		Close() error
	}

	client struct {
		conns
		RedisType RedisType
	}
)

func NewClient(options interface{}) (c Client, err error) {
	switch options := options.(type) {
	case NodeOptions:
		t := &client{}
		t.conns = CreatePool(options.Address, options.OptionsToOption())
		t.RedisType = RedisTypeNode
		_, err = t.do(context.Background(), "PING")
		c = t
	case ClusterOptions:
		t := &client{}
		t.conns, err = NewCluster(options.StartupNodes, options.OptionsToOption())
		t.RedisType = RedisTypeCluster
		c = t
	default:
		err = errors.Errorf("unsupported options(%v)", options)
	}
	return c, err
}

func (c *client) GetConn() Conn {
	return &conn{
		Conn: c.conns.Get(),
	}
}

func (c *client) DoWithContext(ctx context.Context, cmd string, args ...interface{}) (reply interface{}, err error) {
	reply, err = c.do(ctx, cmd, args...)
	return reply, err
}

func (c *client) Close() {
	_ = c.conns.Close()
}

func (c *client) do(ctx context.Context, cmd string, args ...interface{}) (reply interface{}, err error) {
	conn := c.GetConn().GetRealConn()
	defer conn.Close()
	return conn.Do(cmd, args...)
}
