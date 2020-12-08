package main

import (
	"context"
	"github.com/ivalue2333/pkg/src/httpx"
	"github.com/ivalue2333/pkg/src/httpx/middles"
)

var (
	ctx = context.Background()
)

func main() {
	options := httpx.Options{
		Name:    "demo",
		Address: ":8081",
	}

	server := httpx.NewServer(options, httpx.WithMiddles(middles.LoggingRequest(), middles.LoggingResponse()))

	server.GoTask(ctx)
}
