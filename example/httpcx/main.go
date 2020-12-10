package main

import (
	"context"
	"fmt"
	"github.com/ivalue2333/pkg/example/stdout"
	"github.com/ivalue2333/pkg/src/httpcx"
	"github.com/ivalue2333/pkg/src/httpx/httpxcode"
	"github.com/ivalue2333/pkg/src/jsonx"
	"net/url"
	"time"
)

var (
	options = httpcx.Options{
		Address: "http://localhost:8081",
		Timeout: 3,
	}

	client = httpcx.NewClient(httpcx.WithAddress(options.Address), httpcx.WithTimeout(options.Timeout*time.Second))

	ctx = context.Background()

	err error
)

func Get() {
	stdout.PrintFunc("GET")

	type data struct {
		Pong string `json:"pong"`
	}

	type resp struct {
		httpxcode.RespBase
		Data data `json:"data"`
	}

	vals := url.Values{"ping": []string{"123"}}
	r := new(resp)
	err = client.Get(ctx, "/ping", vals, nil, 3, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonx.MarshalUnsafe(r)))
}

func Error() {
	stdout.PrintFunc("Error")

	type resp struct {
		httpxcode.RespBase
	}

	vals := url.Values{"ping": []string{"123"}}
	r := new(resp)
	err = client.Get(ctx, "/error", vals, nil, 3, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonx.MarshalUnsafe(r)))
}

func Post() {
	stdout.PrintFunc("Post")
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

	type resp struct {
		httpxcode.RespBase
		Data PostResp `json:"data"`
	}

	r := new(resp)
	req := &PostReq{Name:"percy", Age: 17}
	err = client.PostJSON(ctx, "/post", req, nil, 3, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonx.MarshalUnsafe(r)))

}

func main() {
	Get()
	Error()
	Post()
}
