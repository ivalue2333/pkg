package es_wrapped

import (
	"context"

	"github.com/olivere/elastic/v7"
)

type WrappedClient struct {
	cc *elastic.Client
}

func NewClient(ctx context.Context, opts ...elastic.ClientOptionFunc) (*WrappedClient, error) {
	//这里不使用自动嗅探否则会转成内网地址或者docker ip导致连不上
	opts = append(opts, elastic.SetSniff(false))
	client, err := elastic.NewClient(opts...)

	if err != nil {
		return nil, err
	}

	return &WrappedClient{cc: client}, nil
}

func (wc *WrappedClient) Ping(ctx context.Context, url string) error {
	_, _, err := wc.cc.Ping(url).Do(ctx)

	return err
}

func (wc *WrappedClient) Client() *elastic.Client { return wc.cc }
