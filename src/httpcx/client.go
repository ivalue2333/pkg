package httpcx

import (
	"context"
	"fmt"
	"gopkg.in/resty.v1"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	GetRestyKernel() (kernel *resty.Client)
	PostJSON(ctx context.Context, path string, values interface{}, headers http.Header, timeoutSeconds int, ret interface{}) (err error)
	Get(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, ret interface{}) (err error)
}

type client struct {
	kernel  *resty.Client
	options Options
}

func NewClient(opts ...Option) Client {
	return NewClientWithOptions(newOptions(opts...))
}

func NewClientWithOptions(options Options) Client {
	var c *client

	// retry to avoid connection refused error, while k8s restart some deploy
	retryTransport := CreateRetryTransport()

	c = &client{
		kernel:  resty.NewWithClient(&http.Client{Transport: retryTransport}),
		options: options,
	}

	c.kernel.SetHostURL(c.options.Address)
	c.kernel.SetTimeout(c.options.Timeout)
	c.kernel.SetRetryCount(c.options.RetryCount)
	c.kernel.SetRetryWaitTime(c.options.RetryWaitTime)
	c.kernel.SetRetryMaxWaitTime(c.options.RetryMaxWaitTime)
	return c
}

func (c *client) GetRestyKernel() *resty.Client {
	return c.kernel
}

func (c *client) PostJSON(ctx context.Context, path string, values interface{}, headers http.Header, timeoutSeconds int, ret interface{}) (err error) {
	if timeoutSeconds > 0 {
		c.kernel.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	}
	r := c.kernel.NewRequest().SetContext(ctx)

	if headers != nil {
		r.Header = headers
	}

	r.SetHeader("Content-Type", "application/json")

	if values != nil {
		r.SetBody(values)
	}

	if ret != nil {
		r.SetResult(ret)
	}

	resp, err := r.Post(path)

	if err != nil {
		return fmt.Errorf("PostJSON:{%s} param:%+v err: %v", r.URL, values, err)
	}
	httpStatusCode := resp.StatusCode()
	if httpStatusCode < http.StatusOK || httpStatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("PostJSON:{%s} param:%+v failed: %v", resp.Request.URL, values, httpStatusCode)
	}

	return
}

func (c *client) Get(ctx context.Context, path string, values url.Values, headers http.Header, timeoutSeconds int, ret interface{}) (err error) {
	if timeoutSeconds > 0 {
		c.kernel.SetTimeout(time.Duration(timeoutSeconds) * time.Second)
	}
	r := c.kernel.NewRequest().SetContext(ctx)

	if headers != nil {
		r.Header = headers
	}

	if values != nil {
		r.QueryParam = values
	}

	if ret != nil {
		r.SetResult(ret)
	}

	resp, err := r.Get(path)
	if err != nil {
		return fmt.Errorf("get:{%s} param:%+v err: %v", r.URL, r.QueryParam, err)
	}
	httpStatusCode := resp.StatusCode()
	if httpStatusCode < http.StatusOK || httpStatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("get:{%s} param:%+v failed: %v", resp.Request.URL, values, httpStatusCode)
	}

	return
}
