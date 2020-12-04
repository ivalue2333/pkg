package mongox

import (
	"context"
	"github.com/ivalue2333/pkg/logx"
	"github.com/ivalue2333/pkg/storex/mongox/mongoc_wrapped"
	"io"
	"time"

	"github.com/tal-tech/go-zero/core/syncx"
	mgOptions "go.mongodb.org/mongo-driver/mongo/options"
)

/*
	not used
 */

const (
	defaultConcurrency = 50
	defaultTimeout     = time.Second
)

var clientManager = syncx.NewResourceManager()

type concurrentClient struct {
	*mongoc_wrapped.WrappedClient
	limit syncx.TimeoutLimit
}

func (cs *concurrentClient) Close() error {
	return nil
}

func getConcurrentClient(url string) (*concurrentClient, error) {
	val, err := clientManager.GetResource(url, func() (io.Closer, error) {
		ctx := context.Background()
		c, err := mongoc_wrapped.NewClient(ctx, mgOptions.Client().ApplyURI(url))
		if err != nil {
			return nil, err
		}

		if err = c.Connect(ctx); err != nil {
			return nil, err
		}

		// ping before use
		if err = c.Ping(ctx, nil); err != nil {
			return nil, err
		}

		cs := &concurrentClient{
			WrappedClient: c,
			limit:         syncx.NewTimeoutLimit(defaultConcurrency),
		}

		return cs, nil

	})
	if err != nil {
		return nil, err
	}

	return val.(*concurrentClient), nil
}

func (cs *concurrentClient) putClient(client *mongoc_wrapped.WrappedClient) {
	if err := cs.limit.Return(); err != nil {
		logx.Error(context.Background(), err)
	}
	return
}

func (cs *concurrentClient) takeClient(opts ...Option) (*mongoc_wrapped.WrappedClient, error) {
	o := &options{
		timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(o)
	}

	if err := cs.limit.Borrow(o.timeout); err != nil {
		return nil, err
	} else {
		return cs.WrappedClient, nil
	}
}
