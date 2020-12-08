package routes

import "github.com/ivalue2333/pkg/src/logx"

type (
	Options struct {
		logger logx.Logger
	}

	Option func(o *Options)
)



