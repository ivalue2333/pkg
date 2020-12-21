package syncx

import "sync"

func OnceFnWithError(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

func OnceFn(fn func()) func() {
	once := new(sync.Once)
	return func() {
		once.Do(fn)
	}
}
