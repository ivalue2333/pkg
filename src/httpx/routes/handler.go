package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ivalue2333/pkg/src/httpx/httpcode"
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

func CreateHandlerFunc(method interface{}, opts ...Option) gin.HandlerFunc {

	mV, reqT, replyT, err := CheckMethod(method)
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

		reply := reflect.New(replyT)

		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req, reply})
		errValue := results[0]
		if errValue.Interface() != nil {
			switch v := errValue.Interface().(type) {
			case error:
				c.JSON(http.StatusOK, RespStruct{
					Code:    httpcode.CodeUnKnown,
					Message: v.Error(),
				})
				return
			default:
				c.JSON(http.StatusOK, RespStruct{
					Code:    httpcode.CodeUnKnown,
					Message: httpcode.CodeUnKnownMessage,
				})
				return
			}
		}

		ret := getRetData(reply)
		c.PureJSON(http.StatusOK, ret)
	}
}

func getRetData(value reflect.Value) interface{} {
	ret := &RespStruct{Code: 0, Data: value.Interface()}
	return ret
}
