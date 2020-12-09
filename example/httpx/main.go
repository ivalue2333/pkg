package main

import (
	"context"
	"errors"
	"github.com/ivalue2333/pkg/src/httpx"
	"github.com/ivalue2333/pkg/src/httpx/httpxroute"
	"github.com/ivalue2333/pkg/src/httpx/middles"
	"net/http"
)

var (
	ctx = context.Background()
)

type PingReq struct {
	Ping string `form:"ping"`
}

type PingResp struct {
	Pong string `json:"pong"`
}

func PingHandler(ctx context.Context, req *PingReq) (*PingResp, error) {
	return &PingResp{Pong: "pong: " + req.Ping}, nil
}

type ErrorReq struct {
}

type ErrorResp struct {
}

func ErrorHandler(ctx2 context.Context, req *ErrorReq) (*ErrorResp, error) {
	return nil, errors.New("error test")
}

func main() {
	options := httpx.Options{
		Name:    "demo",
		Address: ":8081",
	}

	server := httpx.NewServer(options, httpx.WithMiddles(middles.LoggingRequest(), middles.LoggingResponse()))

	engine := server.GetKernel()

	httpxroute.Route(engine, http.MethodGet, "/ping", PingHandler)
	httpxroute.Route(engine, http.MethodGet, "/error", ErrorHandler)

	server.GoTask(ctx)

}
