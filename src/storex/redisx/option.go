package redisx

type Option func(*Options)

func WithConnectTimeout(timeout int) Option {
	return func(opt *Options) {
		opt.ConnectTimeout = timeout
	}
}

func WithReadTimeout(timeout int) Option {
	return func(opt *Options) {
		opt.ReadTimeout = timeout
	}
}

func WithWriteTimeout(timeout int) Option {
	return func(opt *Options) {
		opt.WriteTimeout = timeout
	}
}

func WithIdleTimeout(timeout int) Option {
	return func(opt *Options) {
		opt.IdleTimeout = timeout
	}
}

func WithWait(w bool) Option {
	return func(opt *Options) {
		opt.Wait = w
	}
}

func WithMaxActive(active int) Option {
	return func(opt *Options) {
		opt.MaxActive = active
	}
}

func WithMaxIdle(idle int) Option {
	return func(opt *Options) {
		opt.MaxIdle = idle
	}
}

func WithDB(db int) Option {
	return func(opt *Options) {
		opt.DB = db
	}
}

func WithPassword(password string) Option {
	return func(opt *Options) {
		opt.Password = password
	}
}
