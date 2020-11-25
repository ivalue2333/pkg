package gpool

type Options struct {
	Name          string `json:"name" mapstructure:"name"`
	GNum          int    `json:"g_num" mapstructure:"g_num"`
	ChannelBuffer int    `json:"channel_buffer" mapstructure:"channel_buffer"`
	InputRetry    int    `json:"input_retry" mapstructure:"input_retry"`
}
