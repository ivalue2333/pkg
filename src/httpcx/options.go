package httpcx

import (
	"time"
)

type Options struct {
	Address          string        `mapstructure:"address" json:"address" toml:"address"`
	Timeout          time.Duration `mapstructure:"timeout" json:"timeout" toml:"timeout"`
	RetryCount       int           `mapstructure:"retry_count" json:"retry_count" toml:"retry_count"`                         // 重试次数
	RetryWaitTime    time.Duration `mapstructure:"retry_wait_time" json:"retry_wait_time" toml:"retry_wait_time"`             // 重试间隔等待时间
	RetryMaxWaitTime time.Duration `mapstructure:"retry_max_wait_time" json:"retry_max_wait_time" toml:"retry_max_wait_time"` // 重试间隔最大等待时间
}

func newOptions(opts ...Option) Options {
	options := Options{
		Address:          "",
		Timeout:          3 * time.Second,
		RetryCount:       3,
		RetryWaitTime:    time.Duration(100) * time.Millisecond,
		RetryMaxWaitTime: time.Duration(2000) * time.Millisecond,
	}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type Option func(*Options)

func WithAddress(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithRetryCount(retryCount int) Option {
	return func(options *Options) {
		options.RetryCount = retryCount
	}
}

func WithRetryWaitTime(retryWaitTime time.Duration) Option {
	return func(options *Options) {
		options.RetryWaitTime = retryWaitTime
	}
}

func WithRetryMaxWaitTime(retryMaxWaitTime time.Duration) Option {
	return func(options *Options) {
		options.RetryMaxWaitTime = retryMaxWaitTime
	}
}
