package cachex

import (
	"context"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/ivalue2333/pkg/src/jsonx"
	"github.com/ivalue2333/pkg/src/logx"
	"github.com/ivalue2333/pkg/src/mathx"
	"github.com/ivalue2333/pkg/src/storex/redisx"
	"github.com/ivalue2333/pkg/src/syncx"
	"time"
)

const (
	notFoundPlaceholder = "#*#"
	// make the expiry unstable to avoid lots of cached items expire at the same time
	// make the unstable expiry to be [0.95, 1.05] * seconds
	expiryDeviation = 0.05
)

// indicates there is no such value associate with the key
var errPlaceholder = errors.New("placeholder")

type (
	Cache interface {
		DelCache(ctx context.Context, keys ...string) error
		GetCache(ctx context.Context, key string, v interface{}) error
		SetCache(ctx context.Context, key string, v interface{}) error
		SetCacheWithExpire(ctx context.Context, key string, v interface{}, expire time.Duration) error
		Take(ctx context.Context, v interface{}, key string, query func(v interface{}) error) error
		TakeWithExpire(ctx context.Context, v interface{}, key string, query func(v interface{}, expire time.Duration) error) error
	}

	cacheNode struct {
		client         redisx.Client
		expiry         time.Duration
		errNotFound    error // custom err, for mongo, sql and et...
		notFoundExpiry time.Duration
		unstableExpiry mathx.Unstable
		barrier        syncx.SharedCalls
	}
)

func NewCache(client redisx.Client, err error, opts ...Option) Cache {
	o := newOptions(opts...)
	return cacheNode{
		client:         client,
		errNotFound:    err,
		expiry:         o.Expiry,
		notFoundExpiry: o.NotFoundExpiry,
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		barrier:        syncx.NewSharedCalls(),
	}
}

func (c cacheNode) DelCache(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	vals := make([]interface{}, len(keys))
	for i, key := range keys {
		vals[i] = key
	}

	if _, err := c.client.DoWithContext(ctx, "DEL", vals...); err != nil {
		logx.Errorf(ctx, "failed to clear cache with keys: %q, error: %v", formatKeys(keys), err)
	}

	return nil
}

func (c cacheNode) GetCache(ctx context.Context, key string, v interface{}) error {
	if err := c.doGetCache(ctx, key, v); err == errPlaceholder {
		return c.errNotFound
	} else {
		return err
	}
}

func (c cacheNode) SetCache(ctx context.Context, key string, v interface{}) error {
	return c.SetCacheWithExpire(ctx, key, v, c.aroundDuration(c.expiry))
}

func (c cacheNode) SetCacheWithExpire(ctx context.Context, key string, v interface{}, expire time.Duration) error {
	data, err := jsonx.Marshal(v)
	if err != nil {
		return err
	}
	_, err = c.client.DoWithContext(ctx, "SET", key, string(data), "EX", int(expire.Seconds()))
	return err
}

func (c cacheNode) Take(ctx context.Context, v interface{}, key string, query func(v interface{}) error) error {
	return c.doTake(ctx, v, key, query, func(v interface{}) error {
		return c.SetCache(ctx, key, v)
	})
}

func (c cacheNode) TakeWithExpire(ctx context.Context, v interface{}, key string,
	query func(v interface{}, expire time.Duration) error) error {
	expire := c.aroundDuration(c.expiry)
	return c.doTake(ctx, v, key, func(v interface{}) error {
		return query(v, expire)
	}, func(v interface{}) error {
		return c.SetCacheWithExpire(ctx, key, v, expire)
	})
}

func (c cacheNode) doGetCache(ctx context.Context, key string, v interface{}) error {
	// redis not found, will return redis.ErrNil
	data, err := redis.String(c.client.DoWithContext(ctx, "GET", key))
	if err != nil {
		return err
	}
	if data == notFoundPlaceholder {
		return errPlaceholder
	}
	return c.processCache(ctx, key, data, v)
}

func (c cacheNode) processCache(ctx context.Context, key string, data string, v interface{}) error {
	err := jsonx.Unmarshal([]byte(data), v)
	if err == nil {
		return nil
	}
	return err
}

func (c cacheNode) aroundDuration(duration time.Duration) time.Duration {
	return c.unstableExpiry.AroundDuration(duration)
}

func (c cacheNode) doTake(ctx context.Context, v interface{}, key string,
	query func(v interface{}) error,
	cacheVal func(v interface{}) error) error {

	val, fresh, err := c.barrier.DoEx(key, func() (interface{}, error) {
		if err := c.doGetCache(ctx, key, v); err != nil {
			if err == errPlaceholder {
				return nil, c.errNotFound
			} else if err != redis.ErrNil {
				// this is because redis is in bad status
				// why we just return the error instead of query from db,
				// because we don't allow the disaster pass to the dbs.
				// fail fast, in case we bring down the dbs.
				return nil, err
			}

			// not found, so query
			if err = query(v); err == c.errNotFound {
				if err = c.setCacheWithNotFound(ctx, key); err != nil {
					logx.Error(ctx, "setCacheWithNotFound ", err)
				}
				return nil, c.errNotFound
			} else if err != nil {
				return nil, err
			}

			// set cache
			if err = cacheVal(v); err != nil {
				logx.Error(ctx, "cacheVal ", err)
			}
		}

		return jsonx.Marshal(v)
	})
	if err != nil {
		return err
	}
	if fresh {
		return nil
	}

	return jsonx.Unmarshal(val.([]byte), v)
}

func (c cacheNode) setCacheWithNotFound(ctx context.Context, key string) error {
	_, err := c.client.DoWithContext(ctx, "SET", key, notFoundPlaceholder, int(c.aroundDuration(c.notFoundExpiry).Seconds()))
	return err
}
