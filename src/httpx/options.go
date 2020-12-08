package httpx

import "github.com/ivalue2333/pkg/src/httpx/middles"

type (
	Options struct {
		Name    string `mapstructure:"name" json:"name" toml:"name"`
		Address string `mapstructure:"address" json:"address" toml:"address"`
		Middles []middles.Middle
	}

	Option func(*Options)
)

var (
	defaultOptions = Options{
		Name: "defaultHttpServer",
		Address: ":8080",
		Middles : []middles.Middle{middles.Recovery()},
	}
)

func WithDefault(options *Options)  {
	if options.Name == "" {
		options.Name = defaultOptions.Name
	}
	if options.Address == "" {
		options.Address = defaultOptions.Address
	}
	options.Middles = append(options.Middles, defaultOptions.Middles...)
}

func WithName(name string) Option {
	return func(options *Options) {
		options.Name = name
	}
}

func WithAddress(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}

func WithMiddles(ms ...middles.Middle) Option {
	return func(options *Options) {
		options.Middles = ms
	}
}