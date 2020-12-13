package elasticx

type Options struct {
	ClientName     string   `mapstructure:"client_name" json:"client_name" toml:"client_name"`
	Addrs          []string `mapstructure:"addrs" json:"addrs" toml:"addrs"`
	RequestTimeout int      `mapstructure:"timeout" json:"timeout" toml:"timeout"`
}
