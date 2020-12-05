package redisx

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"time"
)

type (
	// Conn represents a connection to a Redis server.
	Conn interface {
		// Close closes the connection.
		Close() error

		// Err returns a non-nil value when the connection is not usable.
		Err() error

		// Do sends a command to the server and returns the received reply.
		Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error)

		// Send writes the command to the client's output buffer.
		Send(ctx context.Context, commandName string, args ...interface{}) error

		// Flush flushes the output buffer to the Redis server.
		Flush(ctx context.Context) error

		// Receive receives a single reply from the Redis server
		Receive(ctx context.Context) (reply interface{}, err error)

		//Bind
		Bind(keys ...string) error

		//get real redis conn
		GetRealConn() redis.Conn

		// get retry Conn
		GetRetryConn(maxAtt int, tryAgainDelay time.Duration) (redis.Conn, error)
	}

	conn struct {
		redis.Conn
	}
)

func (c *conn) Close() error {
	return c.Conn.Close()
}

func (c *conn) Err() error {
	return c.Conn.Err()
}

func (c *conn) Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.Conn.Do(commandName, args...)
}

func (c *conn) Send(ctx context.Context, commandName string, args ...interface{}) (err error) {
	return c.Conn.Send(commandName, args...)
}

func (c *conn) Flush(ctx context.Context) error {
	return c.Conn.Flush()
}

func (c *conn) Receive(ctx context.Context) (reply interface{}, err error) {
	return c.Conn.Receive()
}

func (c *conn) Bind(keys ...string) error {
	return redisc.BindConn(c.Conn, keys...)
}

func (c *conn) GetRealConn() redis.Conn {
	return c.Conn
}

func (c *conn) GetRetryConn(maxAtt int, tryAgainDelay time.Duration) (redis.Conn, error) {
	return redisc.RetryConn(c.GetRealConn(), maxAtt, tryAgainDelay)
}
