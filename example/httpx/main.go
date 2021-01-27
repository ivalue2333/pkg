package main

import (
	"context"
	"errors"
	"github.com/ivalue2333/pkg/src/httpx"
	"github.com/ivalue2333/pkg/src/httpx/httpxroute"
	"github.com/ivalue2333/pkg/src/httpx/middles"
	"net/http"
	"time"
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

func ErrorHandler(ctx context.Context, req *ErrorReq) (*ErrorResp, error) {
	return nil, errors.New("error test")
}

type PostReq struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sort *int   `json:"sort"`
}

type PostResp struct {
	Name string    `json:"name"`
	Age  int       `json:"age"`
	Sort *int      `json:"sort"`
	Time time.Time `json:"time"`
}

func PostHandler(ctx context.Context, req *PostReq) (*PostResp, error) {
	resp := &PostResp{
		Name: req.Name,
		Age:  req.Age,
		Sort: req.Sort,
		Time: time.Now(),
	}
	return resp, nil
}

func main() {
	var ctx = context.Background()
	options := httpx.Options{
		Name:    "demo",
		Address: ":8081",
	}

	server := httpx.NewServer(httpx.WithName(options.Name), httpx.WithAddress(options.Address),
		httpx.WithMiddles(httpxmiddles.LoggingRequest(), httpxmiddles.LoggingResponse()))

	engine := server.GetKernel()

	httpxroute.Route(engine, http.MethodGet, "/ping", PingHandler)
	httpxroute.Route(engine, http.MethodGet, "/error", ErrorHandler)
	httpxroute.Route(engine, http.MethodPost, "/post", PostHandler)

	server.GoTask(ctx)

}
