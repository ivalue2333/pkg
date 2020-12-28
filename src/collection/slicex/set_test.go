package slicex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListStringContain(t *testing.T) {
	datas := []string{"a", "b", "c"}
	assert.True(t, ListStringContain(datas, "a"))
	assert.False(t, ListStringContain(datas, "d"))
}

func TestListIntContain(t *testing.T) {
	datas := []int{1,2,3}
	assert.True(t, ListIntContain(datas, 1))
	assert.False(t, ListIntContain(datas, 4))
}

func TestListInt64Contain(t *testing.T) {
	datas := []int64{1,2,3}
	assert.True(t, ListInt64Contain(datas, 1))
	assert.False(t, ListInt64Contain(datas, 4))
}
