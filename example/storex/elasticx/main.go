package main

import (
	"context"
	"github.com/ivalue2333/pkg/example/stdout"
	"github.com/ivalue2333/pkg/src/jsonx"
	"github.com/ivalue2333/pkg/src/storex/elasticx"
)

var (
	ctx        = context.Background()
	clientName = "myclient"
	model      *elasticx.Base
)

func init() {
	err := elasticx.ClientsMgr().NewClient(ctx, clientName, []string{"127.0.0.1:9200"}...)
	if err != nil {
		panic(err)
	}
	model = elasticx.NewBaseModelV7(clientName, "myindex")
}

type Data struct {
	Name   string
	Age    int
	School string
}

func insert() {
	stdout.PrintFunc("insert")
	data := &Data{Name: "percy", Age: 12, School: "swjtu"}
	model.Insert(ctx, "", string(jsonx.MarshalUnsafe(data)))
}

func main() {
	insert()
}
