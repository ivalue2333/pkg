package threading

import (
	"bytes"
	"github.com/tal-tech/go-zero/core/rescue"
	"runtime"
	"strconv"
)

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer rescue.Recover()
	fn()
}

// Only for debug, never use it in production
func RoutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	// if error, just return 0
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}