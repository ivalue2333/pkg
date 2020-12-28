package slicex

import "github.com/ivalue2333/pkg/src/collection/setx"

// string contain
func ListStringContain(list []string, val string) bool {
	for _, data := range list {
		if data == val {
			return true
		}
	}
	return false
}

// int contain
func ListIntContain(list []int, val int) bool {
	for _, data := range list {
		if data == val {
			return true
		}
	}
	return false
}

// int64 contain
func ListInt64Contain(list []int64, val int64) bool {
	for _, data := range list {
		if data == val {
			return true
		}
	}
	return false
}

// list string to set
func ListStringSet(list []string) *setx.ItemSet {
	set := new(setx.ItemSet)
	for _, data := range list {
		set.Add(data)
	}
	return set
}

// list int to set
func ListIntSet(list []int) *setx.ItemSet {
	set := new(setx.ItemSet)
	for _, data := range list {
		set.Add(data)
	}
	return set
}

// list int64 to set
func ListInt64Set(list []int) *setx.ItemSet {
	set := new(setx.ItemSet)
	for _, data := range list {
		set.Add(data)
	}
	return set
}


