package slicex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange(t *testing.T) {
	sl1 := Range(0, 1)
	assert.Equal(t, len(sl1), 1)

	N := 100
	sl2 := Range(0, N)
	assert.Equal(t, len(sl2), N)
	for i := 0; i < N; i++ {
		assert.Equal(t, sl2[i], i)
	}

	start := 10
	end := 20
	sl3 := Range(start, end)
	assert.Equal(t, len(sl3), end-start)
	for i := 0; i < end-start; i++ {
		assert.Equal(t, sl3[i], i+start)
	}
}

func TestRangeInternal(t *testing.T) {
	tests := []struct {
		start    int
		end      int
		internal int
		length   int
		max      int
		min      int
	}{
		{start: 0, end: 100, internal: 2, length: 50, max: 98, min: 0},
		{start: 0, end: 100, internal: 3, length: 34, max: 99, min: 0},
		{start: 10, end: 110, internal: 2, length: 50, max: 108, min: 10},
		{start: 10, end: 110, internal: 3, length: 34, max: 109, min: 10},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			sls := RangeInternal(test.start, test.end, test.internal)
			//t.Log(sls)
			assert.Equal(t, test.max, sls[len(sls)-1])
			assert.Equal(t, test.min, sls[0])
			assert.Equal(t, test.length, len(sls))
		})
	}
}
