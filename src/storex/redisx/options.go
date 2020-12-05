package redisx

type Options struct {
	ConnectTimeout int    `json:"connect_timeout" toml:"connect_timeout" example:"2"`
	ReadTimeout    int    `json:"read_timeout" toml:"read_timeout" example:"2"`
	WriteTimeout   int    `json:"write_timeout" toml:"write_timeout" example:"2"`
	MaxActive      int    `json:"max_active" toml:"max_active" example:"200"`
	MaxIdle        int    `json:"max_idle" toml:"max_idle" example:"200"`
	IdleTimeout    int    `json:"idle_timeout" toml:"idle_timeout" example:"2"`
	Wait           bool   `json:"wait" toml:"wait" example:"false"`
	DB             int    `json:"db" toml:"db"`
	Password       string `json:"password" toml:"password"`
}

// redis single node
type NodeOptions struct {
	Address        string `mapstructure:"address" json:"address" toml:"address"`
	DB             int    `mapstructure:"db" json:"db" toml:"db"`
	Password       string `mapstructure:"password" json:"password" toml:"password"`
	ConnectTimeout int    `mapstructure:"connect_timeout" json:"connect_timeout" toml:"connect_timeout"`
	ReadTimeout    int    `mapstructure:"read_timeout" json:"read_timeout" toml:"read_timeout"`
	WriteTimeout   int    `mapstructure:"write_timeout" json:"write_timeout" toml:"write_timeout"`
	MaxActive      int    `mapstructure:"max_active" json:"max_active" toml:"max_active"`
	MaxIdle        int    `mapstructure:"max_idle" json:"max_idle" toml:"max_idle"`
	IdleTimeout    int    `mapstructure:"idle_timeout" json:"idle_timeout" toml:"idle_timeout"`
	Wait           bool   `mapstructure:"wait" json:"wait" toml:"wait"`
}

// redis cluster
type ClusterOptions struct {
	StartupNodes   []string `mapstructure:"startup_nodes" json:"startup_nodes" toml:"startup_nodes"`
	DB             int      `mapstructure:"db" json:"db" toml:"db"`
	Password       string   `mapstructure:"password" json:"password" toml:"password"`
	ConnectTimeout int      `mapstructure:"connect_timeout" json:"connect_timeout" toml:"connect_timeout"`
	ReadTimeout    int      `mapstructure:"read_timeout" json:"read_timeout" toml:"read_timeout"`
	WriteTimeout   int      `mapstructure:"write_timeout" json:"write_timeout" toml:"write_timeout"`
	MaxActive      int      `mapstructure:"max_active" json:"max_active" toml:"max_active"`
	MaxIdle        int      `mapstructure:"max_idle" json:"max_idle" toml:"max_idle"`
	IdleTimeout    int      `mapstructure:"idle_timeout" json:"idle_timeout" toml:"idle_timeout"`
	Wait           bool     `mapstructure:"wait" json:"wait" toml:"wait"`
}

func (nodeOptions NodeOptions) OptionsToOption() Option {
	return func(option *Options) {
		if nodeOptions.ConnectTimeout != 0 {
			option.ConnectTimeout = nodeOptions.ConnectTimeout
		}
		if nodeOptions.ReadTimeout != 0 {
			option.ReadTimeout = nodeOptions.ReadTimeout
		}
		if nodeOptions.WriteTimeout != 0 {
			option.WriteTimeout = nodeOptions.WriteTimeout
		}
		if nodeOptions.MaxActive != 0 {
			option.MaxActive = nodeOptions.MaxActive
		}
		if nodeOptions.MaxIdle != 0 {
			option.MaxIdle = nodeOptions.MaxIdle
		}
		if nodeOptions.IdleTimeout != 0 {
			option.IdleTimeout = nodeOptions.IdleTimeout
		}
		if nodeOptions.DB != 0 {
			option.DB = nodeOptions.DB
		}
		if nodeOptions.Password != "" {
			option.Password = nodeOptions.Password
		}
		option.Wait = nodeOptions.Wait
	}
}

func (clusterOptions ClusterOptions) OptionsToOption() Option {
	return func(option *Options) {
		if clusterOptions.ConnectTimeout != 0 {
			option.ConnectTimeout = clusterOptions.ConnectTimeout
		}
		if clusterOptions.ReadTimeout != 0 {
			option.ReadTimeout = clusterOptions.ReadTimeout
		}
		if clusterOptions.WriteTimeout != 0 {
			option.WriteTimeout = clusterOptions.WriteTimeout
		}
		if clusterOptions.MaxActive != 0 {
			option.MaxActive = clusterOptions.MaxActive
		}
		if clusterOptions.MaxIdle != 0 {
			option.MaxIdle = clusterOptions.MaxIdle
		}
		if clusterOptions.IdleTimeout != 0 {
			option.IdleTimeout = clusterOptions.IdleTimeout
		}
		if clusterOptions.Password != "" {
			option.Password = clusterOptions.Password
		}
		option.Wait = clusterOptions.Wait
	}
}
