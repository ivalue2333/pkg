package main

import (
	"fmt"
	"github.com/ivalue2333/pkg/src/collection/slicex"
	"github.com/ivalue2333/pkg/src/mr"
	"log"
	"sync/atomic"
)

func MapReduceDemoSum() {
	ints := slicex.Range(0, 100)
	var total int64

	generator := func(source chan<- interface{}) {
		for _, i := range ints {
			source <- i
		}
	}

	mapperFn := func(item interface{}) {
		i := item.(int)
		atomic.AddInt64(&total, int64(i))
	}

	mr.MapVoid(generator, mapperFn)

	fmt.Println(total)
}

func MapReduceDemoCheckLegal() ([]int, error) {

	uids := slicex.Range(0, 1000)

	r, err := mr.MapReduce(func(source chan<- interface{}) {
		for _, uid := range uids {
			source <- uid
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		uid := item.(int)
		ok, err := check(uid)
		if err != nil {
			cancel(err)
		}
		if ok {
			writer.Write(uid)
		}
	}, func(pipe <-chan interface{}, writer mr.Writer, cancel func(error)) {
		var uids []int
		for p := range pipe {
			uids = append(uids, p.(int))
		}
		writer.Write(uids)
	})
	if err != nil {
		log.Printf("check error: %v", err)
		return nil, err
	}

	return r.([]int), nil
}

func check(uid int) (bool, error) {
	// do something check user legal
	return true, nil
}

func main() {
	//ans, err := MapReduceDemoCheckLegal()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(ans)


	MapReduceDemoSum()
}
