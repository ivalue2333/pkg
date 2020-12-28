package jsonx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalUnsafe(t *testing.T) {
	data := "{abc}"
	fmt.Println(MarshalUnsafeString(data))
	assert.Equal(t, MarshalUnsafe(data), []byte("\"{abc}\""))
}

func TestMarshalUnsafe2(t *testing.T) {
	data := make(map[interface{}]interface{})
	data["k1"] = "ddd"
	data["k2"] = "ddd"
	assert.Equal(t, MarshalUnsafe(data), []byte{})
}

// 如果源中的值与目标中的值不对应，则JSON解码器不会报告错误。
func TestUnmarshal(t *testing.T) {
	type A struct {
		name string `json:"name"`
	}
	var jsonString string = `{"status":false}`
	var a A
	err := Unmarshal([]byte(jsonString), &a)
	assert.Nil(t, err)
}

// 语法错误 报错
func TestUnmarshal2(t *testing.T) {
	type A struct {
		Name string `json:"name"`
	}
	data := []byte(`{"name":what?}`)
	var a A
	err := Unmarshal(data, &a)
	assert.NotNil(t, err)
	t.Log(err.Error())
}

// 类型不匹配
func TestUnmarshal3(t *testing.T) {
	data := []byte(`{"name":false}`)
	type B struct {
		Name string `json:"name"`
	}
	var b B
	err := Unmarshal(data, &b)
	assert.NotNil(t, err)
	t.Log(err.Error())
}
