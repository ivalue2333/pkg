package logx

type Config struct {
	Options Options `mapstructure:",squash"`
}

type Options struct {
	Level    string `mapstructure:"level" json:"level" toml:"level"`
	File     string `mapstructure:"file" json:"file" toml:"file"`
	ErrFile  string `mapstructure:"err_file" json:"err_file" toml:"err_file"`
	AppName  string `mapstructure:"app_name" json:"app_name" toml:"app_name"`
	Format   string `mapstructure:"format" json:"format" toml:"format"`
	AutoInit bool   `mapstructure:"auto_init" json:"auto_init" toml:"auto_init"`
}

func newOptions(opts ...Option) Options {
	options := Options{
		Level:   "",
		File:    "",
		ErrFile: "",
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

type Option func(*Options)

func WithLevel(level string) Option {
	return func(options *Options) {
		options.Level = level
	}
}

func WithFile(file string) Option {
	return func(options *Options) {
		options.File = file
	}
}

func WithErrFile(errFile string) Option {
	return func(options *Options) {
		options.ErrFile = errFile
	}
}
