package httpxroute

import "github.com/ivalue2333/pkg/src/logx"

type (
	HandlerOptions struct {
		logger       logx.Logger
		costumedCode bool
	}

	HandlerOption func(o *HandlerOptions)
)


func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{}
}

func WithCostumedCode(costumedCode bool) HandlerOption {
	return func(o *HandlerOptions) {
		o.costumedCode = costumedCode
	}
}
