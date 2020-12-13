package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ivalue2333/pkg/src/jsonx"
	"github.com/ivalue2333/pkg/src/storex/cachex"
	"github.com/ivalue2333/pkg/src/storex/redisx"
	"os"
)

var (
	redisc  redisx.Client
	ctx     = context.Background()
	ErrMock = errors.New("Mock: nil return")

	KeyMock  = "{apartment}:{service}percy"
	DataMock = Data{Name: "percy001", Age: 10}
)

func init() {
	option := redisx.NodeOptions{
		//Address: "127.0.0.1:6379",
		Address: os.Getenv("REDIS_URI"),
	}
	var err error
	redisc, err = redisx.NewClient(option)
	if err != nil {
		panic(err)
	}
}

type (
	Data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
)

func get() {
	cachec := cachex.NewCache(redisc, ErrMock)
	data := &Data{}
	var err error
	err = cachec.GetCache(ctx, KeyMock, data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func take() {
	cachec := cachex.NewCache(redisc, ErrMock)
	data := &Data{}
	var err error

	query := func(v interface{}) error {
		fmt.Println("hits")
		bData := jsonx.MarshalUnsafe(DataMock)
		return jsonx.Unmarshal(bData, v)
	}

	err = cachec.Take(ctx, data, KeyMock, query)
	if err != nil {
		panic(err)
	}

	err = cachec.Take(ctx, data, KeyMock, query)
	if err != nil {
		panic(err)
	}
}

func main() {
	get()
	//take()
}
