package threading

import (
	"github.com/ivalue2333/pkg/src/collection/lang"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSafe(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan lang.PlaceholderType)
	go RunSafe(func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestRoutineId(t *testing.T) {
	assert.True(t, RoutineId() > 0)
}
