package gpool

import (
	"context"
	"testing"
	"time"
)

func TestWorker_RunPanic(t *testing.T) {
	count := 0
	worker := NewWorker(func(data interface{}) {
		count += 1
		if data == "1" {
			panic("panic")
		}
	})
	ctx := context.TODO()
	go worker.Run(ctx)
	datas := []string{"1", "2", "3", "4", "5"}
	for _, data := range datas {
		worker.Input(ctx, data)
	}

	time.Sleep(2 * time.Second)

	if count != len(datas) {
		t.Errorf("failed: count:%d, len:%d", count, len(datas))
	}
	t.Log("done")
}

func TestWorker_RunPanic2(t *testing.T) {
	count := 0
	worker := NewWorker(func(data interface{}) {
		count += 1
		panic(count)
	})
	data := 50
	ctx := context.TODO()
	go worker.Run(ctx)
	for i := 0; i < data; i++ {
		worker.Input(ctx, i)
	}

	time.Sleep(2 * time.Second)

	if count != data {
		t.Errorf("failed: count:%d, len:%d", count, data)
	}
	t.Log("done")
}
