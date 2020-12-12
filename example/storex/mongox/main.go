package main

import (
	"context"
	"fmt"
	"github.com/ivalue2333/pkg/example/stdout"
	"github.com/ivalue2333/pkg/src/storex/mongox"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
	"time"
)

type Template struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Num        float64            `bson:"num" json:"num"`
	Name       string             `bson:"name" json:"name"`
	Day        int                `json:"day" bson:"day"`
	Date       *time.Time         `json:"date" bson:"date"`
	TagList    []string           `json:"tag_list" bson:"tag_list"`
	CreateTime time.Time          `json:"create_time" bson:"create_time"`
}

var (
	model *mongox.Model
	ctx   = context.Background()
)

func init() {
	url := os.Getenv("MONGO_URI")
	ctx := context.Background()
	model = mongox.MustNewModel(ctx, url, "mydb", "mytable")
}

func insert(n int, day int) {
	stdout.PrintFunc("insert")
	for i := 0; i < n; i++ {
		for j := 0; j < day; j++ {
			template := Template{}
			template.Id = primitive.NewObjectID()
			template.Num = float64(i)
			template.Name = "percy" + fmt.Sprintf("%d", i)
			template.Day = j
			template.CreateTime = time.Now()
			if res, err := model.InsertOne(ctx, template); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(*res)
			}
		}
	}
}

func queryOne() {
	stdout.PrintFunc("queryOne")
	res := &Template{}
	err := model.FindOne(ctx, primitive.M{}, res)
	fmt.Println(res)
	fmt.Println(err == mongo.ErrNoDocuments)

	err = model.FindOne(ctx, primitive.M{"name": "alice"}, res)
	fmt.Println(err == mongo.ErrNoDocuments)
}

func queryMany() {
	stdout.PrintFunc("queryMany")

	res := make([]Template, 0)
	opts := options.Find()
	skip, limit := int64(0), int64(3)
	opts.Skip = &skip
	opts.Limit = &limit
	// 使用 1, -1 來实现 顺序排序和倒序排序
	opts.Sort = primitive.M{"num": -1}
	err := model.Find(ctx, primitive.M{}, &res, opts)
	fmt.Println(res)
	fmt.Println(err)
}

func queryManyList() {
	stdout.PrintFunc("queryManyList")

	res := make([]Template, 0)
	err := model.Find(ctx, primitive.M{}, &res)
	fmt.Println(res)
	fmt.Println(err)
}

func aggregate() {
	stdout.PrintFunc("aggregate")

	type TT struct {
		Id         int         `json:"id" bson:"_id"`
		Num        float64     `bson:"num" json:"num"`
		Nums       []float64   `json:"nums" bson:"nums"`
		CreateTime []time.Time `json:"create_time" bson:"create_time"`
	}

	res := make([]TT, 0)
	pipes := []primitive.M{
		primitive.M{"$match": primitive.M{"name": "percy1"}},
		primitive.M{"$group": primitive.M{
			"_id":         "$day",
			"create_time": primitive.M{"$push": "$create_time"},
			"nums":        primitive.M{"$push": "$num"},
			"num":         primitive.M{"$sum": "$num"}},
		},
	}
	err := model.Aggregate(ctx, pipes, &res)
	fmt.Println(len(res), res)
	fmt.Println(err)
}

func aggregate2() {
	stdout.PrintFunc("aggregate2")

	type Idt struct {
		Name string `bson:"name" json:"name"`
		Day  int    `json:"day" bson:"day"`
	}

	// 注意这个结构，除了 Id，其他字段不要是结构体
	type TT struct {
		Id         Idt         `json:"id" bson:"_id"`
		Num        float64     `bson:"num" json:"num"`
		Nums       []float64   `json:"nums" bson:"nums"`
		CreateTime []time.Time `json:"create_time" bson:"create_time"`
	}

	res := make([]TT, 0)
	pipes := []primitive.M{
		primitive.M{"$match": primitive.M{"name": primitive.M{"$in": []string{"percy1", "percy0"}}}},
		primitive.M{"$group": primitive.M{
			"_id":         primitive.M{"name": "$name", "day": "$day"},
			"create_time": primitive.M{"$push": "$create_time"},
			"nums":        primitive.M{"$push": "$num"},
			"num":         primitive.M{"$sum": "$num"}},
		},
	}
	err := model.Aggregate(ctx, pipes, &res)
	fmt.Println(len(res), res)
	fmt.Println(err)
}

func upsertBad() {
	stdout.PrintFunc("upsert")

	t := &Template{}
	t.Name = "percy1234"
	t.Num = 12

	res, err := model.Upsert(ctx, primitive.M{"name": t.Name}, primitive.M{"$set": t})
	fmt.Println(res, err)
}

func upsertGood() {
	stdout.PrintFunc("upsert")

	t := &Template{}
	t.Id = primitive.NewObjectID()
	t.Name = "percy1234"
	t.Num = 12

	res, err := model.Upsert(ctx, primitive.M{"_id": t.Id}, primitive.M{"$set": t})
	fmt.Println(res, err)
}

func updateOne() {
	stdout.PrintFunc("updateOne")

	t := &Template{}
	t.Id = primitive.NewObjectID()
	t.Name = "percy"
	t.Num = 12
	res, err := model.UpdateOne(ctx, primitive.M{"_id": t.Id}, primitive.M{"$set": t})
	fmt.Println(res, err)
}

func deleteMany() {
	stdout.PrintFunc("deleteMany")
	res, err := model.DeleteMany(ctx, primitive.M{})
	fmt.Println(res, err)
}

func shardingInsert() {
	stdout.PrintFunc("shardingInsert")
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		t := &Template{Id: primitive.NewObjectID(), Name: "percy_" + strconv.Itoa(i)}
		_, err := model.InsertOne(ctx, t)
		if err != nil {
			fmt.Println(err, "InsertOne")
		}
	}
}

func shardingFind() {
	stdout.PrintFunc("shardingInsert")
	ctx := context.Background()
	t1 := &Template{}
	var err error
	err = model.FindOne(ctx, primitive.M{"name": "percy_59"}, t1)
	if err != nil {
		fmt.Println(err, "FindOne")
		return
	}
	fmt.Println("t1", *t1)

	t2 := &Template{}
	err = model.FindOne(ctx, primitive.M{}, t2)
	if err != nil {
		fmt.Println(err, "FindOne")
		return
	}
	fmt.Println("t2", *t2)
}

func main() {
	//insert(3, 3)
	//queryOne()
	//queryMany()
	//aggregate()
	//upsertBad()
	//upsertGood()
	//updateOne()
	deleteMany()
	//shardingInsert()
	//shardingFind()
}
