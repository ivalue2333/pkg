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
	"time"
)

var (
	ctx        = context.Background()
	clientName = "myclient"
	indexName  = "myindex"
	model      *elasticx.Base
)

func init() {
	err := elasticx.ClientsMgr().NewClient(ctx, clientName, []string{"http://127.0.0.1:9200"}...)
	if err != nil {
		panic(err)
	}
	model = elasticx.NewBaseModelV7(clientName, indexName)
}

type Data struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	School      string `json:"school"`
	Father      string `json:"father"`
	Email       string `json:"email"`
	Random      string `json:"random"`
	HasChildren bool   `json:"has_children"`
	Article     string `json:"article"`
}

func (d *Data) key() string {
	return d.Name
}

var NAME_DEFAULT = []string{"percy", "alice", "bob", "percy1", "alice1", "bob1"}
var SCHOOL_LIST = []string{"swjtu", "bd", "qinghua", "chuanda", "chengduda", "shifangda"}

func mock() {
	length := 100
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

func deleteIndex() {
	stdout.PrintFunc("deleteIndex")
	res, err := model.Client().DeleteIndex(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func indexSettings() {
	stdout.PrintFunc("indexSettings")
	settings := `{"mappings": {"properties": {"email": {"type": "keyword"}, "name": {"type": "text"}, "age": {"type": "integer"}}}}`
	res, err := model.Client().CreateIndex(indexName).Body(settings).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func indexGet() {
	stdout.PrintFunc("indexGet")
	res, err := model.Client().IndexGet(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonx.MarshalUnsafeString(res[indexName]))
}

func indexGetSettings() {
	stdout.PrintFunc("indexGetSettings")
	res, err := model.Client().IndexGetSettings(indexName).Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonx.MarshalUnsafeString(res[indexName]))
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

func del() {
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
	esParam.PageSize = 2

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

func scroll() {
	stdout.PrintFunc("scroll")
	boolQ := elastic.NewBoolQuery()
	defer func() {
		err := model.CloseScroll(ctx)
		if err != nil {
			fmt.Println("CloseScroll", err)
		}
	}()
	var results []Data
	total, scrollId, err := model.Scroll(ctx, "1m", []string{"name"}, boolQ, map[string]interface{}{"size": 5}, &results)
	if err != nil {
		panic(err)
	}
	fmt.Println(total, len(results), results)
	fmt.Println(scrollId)
}

func scrollAll() {
	stdout.PrintFunc("scrollAll")
	defer func() {
		err := model.CloseScroll(ctx)
		fmt.Println("CloseScroll done")
		if err != nil {
			fmt.Println("scrollAll", err)
		}
	}()
	var results []Data
	var scrollId string
	for {
		scrollId, datas := scrollAll_(scrollId)
		fmt.Println(fmt.Sprintf("length(%d), scrollId(%s)", len(datas), scrollId))
		results = append(results, datas...)
		if scrollId == "" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("results:", len(results))
}

func scrollAll_(scrollId string) (newScrollId string, results []Data) {
	boolQ := elastic.NewBoolQuery()
	if scrollId != "" {
		model.ScrollService().ScrollId(scrollId)
	}
	_, newScrollId, err := model.Scroll(ctx, "1m", []string{"name"}, boolQ, map[string]interface{}{"size": 100}, &results)
	if err != nil {
		//panic(err)
	}
	return newScrollId, results
}

func main() {

	// index
	//deleteIndex()
	//indexSettings()
	//indexGet()
	//indexGetSettings()


	//mock()

	//insert()
	//insertBodyJSON()
	//get()
	//update()
	//upsert()
	//del()
	//search()
	//scroll()
	scrollAll()
}
