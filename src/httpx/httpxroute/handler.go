package httpxroute

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ivalue2333/pkg/src/httpx/httpxcode"
	"net/http"
	"reflect"
)

type requestKey struct{}
type responseKey struct{}
type Initer interface {
	Init(ctx context.Context)
}

var (
	RequestKey  = requestKey{}
	ResponseKey = responseKey{}
)

func CreateHandlerFunc(method interface{}, opts ...HandlerOption) gin.HandlerFunc {

	options := NewHandlerOptions()
	for _, opt := range opts {
		opt(options)
	}

	mV, reqT, _, err := CheckMethod(method)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		req := reflect.New(reqT)
		if err := c.ShouldBind(req.Interface()); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		ctx = context.WithValue(ctx, RequestKey, c.Request)
		ctx = context.WithValue(ctx, ResponseKey, c.Writer)

		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req})
		errValue := results[1]
		if errValue.Interface() != nil {
			switch v := errValue.Interface().(type) {
			case error:
				c.JSON(http.StatusOK, httpxcode.RespBase{
					Code:    httpxcode.CodeBusy,
					Message: v.Error(),
				})
				return
			default:
				c.JSON(http.StatusOK, httpxcode.RespBase{
					Code:    httpxcode.CodeBusy,
					Message: httpxcode.CodeBusyMessage,
				})
				return
			}
		}
		reply := results[0]

		ret := getRetData(reply, options)
		c.PureJSON(http.StatusOK, ret)
	}
}

func getRetData(value reflect.Value, options *HandlerOptions) interface{} {
	if options.costumedCode {
		return value
	}
	base := httpxcode.NewRespBase(httpxcode.CodeOK, httpxcode.CodeOkMessage)
	ret := &httpxcode.RespStruct{Data: value.Interface()}
	ret.RespBase = base
	return ret
}
