package main

import (
	"context"
	"fmt"
	"github.com/ivalue2333/pkg/example/stdout"
	"github.com/ivalue2333/pkg/src/jsonx"
	"github.com/ivalue2333/pkg/src/storex/elasticx"
	"github.com/ivalue2333/pkg/src/stringx"
	"github.com/olivere/elastic/v7"
	"math/rand"
)

var (
	ctx        = context.Background()
	clientName = "myclient"
	model      *elasticx.Base
)

func init() {
	err := elasticx.ClientsMgr().NewClient(ctx, clientName, []string{"http://127.0.0.1:9200"}...)
	if err != nil {
		panic(err)
	}
	model = elasticx.NewBaseModelV7(clientName, "myindex")
}

type Data struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	School      string `json:"school"`
	Father      string `json:"father"`
	Random      string `json:"random"`
	HasChildren bool   `json:"has_children"`
}

func (d *Data) key() string {
	return d.Name
}

var NAME_DEFAULT = []string{"percy", "alice", "bob", "tom"}
var SCHOOL_LIST = []string{"swjtu", "bd", "qinghua", "chuanda", "chengduda", "shifangda"}

func mock() {
	length := rand.Intn(100) + 100
	for i := 0; i < length; i++ {
		for _, name := range NAME_DEFAULT {
			for _, school := range SCHOOL_LIST {
				data := &Data{Name: name, Age: rand.Intn(50), School: school, Random: stringx.Rand()}
				err := model.Insert(ctx, data.key(), string(jsonx.MarshalUnsafe(data)))
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func insert() {
	stdout.PrintFunc("insert")
	for _, name := range NAME_DEFAULT {
		data := &Data{Name: name, Age: 12, School: "swjtu", Random: stringx.Rand()}
		err := model.Insert(ctx, data.key(), string(jsonx.MarshalUnsafe(data)))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func insertBodyJSON() {
	stdout.PrintFunc("insertBodyJSON")
	data := &Data{Name: "insertBodyJSON", Age: 15, School: "swjtu", Random: stringx.Rand()}
	err := model.InsertBodyJSON(ctx, data.key(), "", data)
	if err != nil {
		fmt.Println(err)
	}
}

func get() {
	stdout.PrintFunc("get")
	data := &Data{}
	err := model.Get(ctx, NAME_DEFAULT[0], "", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonx.MarshalUnsafeString(data))
}

func update() {
	stdout.PrintFunc("update")
	data := &Data{Name: NAME_DEFAULT[1]}
	err := model.Update(ctx, data.key(), map[string]interface{}{"age": 18})
	if err != nil {
		fmt.Println(err)
	}
}

func upsert() {
	stdout.PrintFunc("upsert")
	data := &Data{Name: "upsert", Age: 21, School: "swjtu", Random: stringx.Rand()}
	err := model.Upsert(ctx, data.key(), data)
	if err != nil {
		fmt.Println(err)
	}

	data = &Data{Name: "upsert", Age: 1000, School: "swjtu", Random: stringx.Rand(), HasChildren: true}
	err = model.Upsert(ctx, data.key(), data)
	if err != nil {
		fmt.Println(err)
	}
}

func delete() {
	stdout.PrintFunc("delete")
	err := model.Delete(ctx, NAME_DEFAULT[0])
	if err != nil {
		fmt.Println(err, NAME_DEFAULT[0])
	}

	err = model.Delete(ctx, "xxx_not_found")
	if err != nil {
		fmt.Println(err, "xxx_not_found")
	}
}

func search() {
	stdout.PrintFunc("search")
	boolQ := elastic.NewBoolQuery()
	boolQ.Filter(elastic.NewTermQuery("school", "swjtu"))

	esParam := &elasticx.QueryParam{}
	// https://my.oschina.net/dabird/blog/1926042
	//esParam.Sort = "name"
	//esParam.Asc = true
	esParam.PageNo = 0
	esParam.PageSize = 20

	var results []Data

	fn := func(h *elastic.SearchHit) {
		var s Data
		err := jsonx.Unmarshal([]byte(h.Source), &s)
		if err != nil {
			fmt.Println("unmarshal spu failed, err :", err)
			return
		}
		results = append(results, s)
	}

	total, err := model.FindWithForeach(ctx, boolQ, esParam, fn)
	if err != nil {
		fmt.Println("FindWithForeach failed", err)
	}
	fmt.Println(total)
	fmt.Println(results)
}

func main() {
	//insert()
	//insertBodyJSON()
	//get()
	//update()
	//upsert()
	//delete()
	search()
}
