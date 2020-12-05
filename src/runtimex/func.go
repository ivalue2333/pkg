package runtimex

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// FunctionName will get function name and return
func FunctionName(i interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	// 去除路径
	fn = filepath.Ext(fn)
	// 去除.
	fn = strings.TrimPrefix(fn, ".")
	return fn
}
