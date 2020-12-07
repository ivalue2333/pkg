package main

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/ivalue2333/pkg/example/stdout"
	"github.com/ivalue2333/pkg/src/storex/redisx"
	"os"
	"reflect"
	"time"
)

var (
	redisc redisx.Client
	ctx    = context.Background()
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

func setAndGet() {
	stdout.PrintFunc("set")
	data, err := redis.String(redisc.DoWithContext(ctx, "SET", "name", "percy"))
	fmt.Println(data, err)

	data, err = redis.String(redisc.DoWithContext(ctx, "GET", "name"))
	fmt.Println(data, err)
}

func setNx() {
	reply, err := redisc.DoWithContext(ctx, "SET", "user", "percy", "ex", 10, "nx")
	fmt.Println(reply, reflect.TypeOf(reply), err)

	reply, err = redisc.DoWithContext(ctx, "SET", "user", "percy", "ex", 10, "nx")
	fmt.Println(reply, reflect.TypeOf(reply), err)

	data, err := redis.String(redisc.DoWithContext(ctx, "SET", "user2", "percy2", "ex", 10, "nx"))
	fmt.Println(data, err)
}

func setEx() {
	var err error
	reply, err := redisc.DoWithContext(ctx, "SET", "password", "123456", "EX", "1")
	fmt.Println(reply, err)

	time.Sleep(2 * time.Second)

	password, err := redis.String(redisc.DoWithContext(ctx, "GET", "password"))
	fmt.Println("password:->", password, ".err:->", err, ".is true:->", err == redis.ErrNil)
}

func setIntVal() {
	reply, err := redisc.DoWithContext(ctx, "SET", "int", 3)
	fmt.Println(reply, reflect.TypeOf(reply), err)

	reply, err = redisc.DoWithContext(ctx, "GET", "int")

	r1, err1 := redis.Int(reply, err)
	fmt.Println(r1, err1)

	r2, err2 := redis.Int64(reply, err)
	fmt.Println(r2, err2)
}

func del() {
	reply, err := redisc.DoWithContext(ctx, "DEL", "user")
	fmt.Println(reply, reflect.TypeOf(reply), err)
	fmt.Println(redis.Int(reply, err))
}

func hash() {
	reply, err := redisc.DoWithContext(ctx, "HMSET", "user001", "name", "jim", "gender", "main", "age", 12)
	fmt.Println(reply, err)

	reply, err = redis.String(redisc.DoWithContext(ctx, "HGET", "user001", "name"))
	fmt.Println(reply, err)

	reply1, err1 := redis.Strings(redisc.DoWithContext(ctx, "HMGET", "user001", "name", "gender"))
	fmt.Println(reply1, err1)
}

func hashScan() {

}

func pipe() {
	conn := redisc.GetConn()
	var err error

	err = conn.Send(ctx, "SET", "user1", "percy1")
	err = conn.Send(ctx, "SET", "user2", "percy2")
	err = conn.Send(ctx, "GET", "user1")
	err = conn.Flush(ctx)
	fmt.Println(err)
	reply, err := conn.Do(ctx, "GET", "user1")
	fmt.Println(reflect.ValueOf(reply), reflect.TypeOf(reply), string(reply.([]byte)))
}

func main() {
	//setAndGet()
	//setNx()
	//setEx()
	//setIntVal()
	//del()
	//pipe()
	//hash()
}
