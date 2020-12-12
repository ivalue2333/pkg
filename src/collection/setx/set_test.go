package setx

import (
	"fmt"
	"math/rand"
	"testing"
)

func populateSet(count int, start int) *ItemSet {
	set := NewItemSet()
	for i := start; i < (start + count); i++ {
		set.Add(fmt.Sprintf("item%d", i))
	}
	return set
}

func addSet(set *ItemSet, count int, start int) *ItemSet {
	for i := start; i < (start + count); i++ {
		set.Add(fmt.Sprintf("item%d", i))
	}
	return set
}

func opSet(set *ItemSet, count int, start int) *ItemSet {
	for i := start; i < (start + count); i++ {
		set.Add(fmt.Sprintf("item%d", i))
	}
	for i := start; i < (start + count); i++ {
		set.Delete(fmt.Sprintf("item%d", i))
	}
	for i := start; i < (start + count); i++ {
		set.Has(fmt.Sprintf("item%d", i))
	}
	set.Items()
	set.Size()
	return set
}

func TestAdd(t *testing.T) {
	set := populateSet(3, 0)
	if size := set.Size(); size != 3 {
		t.Errorf("wrong count, expected 3 and got %d", size)
	}
	set.Add("item1") //should not add it, already there
	if size := set.Size(); size != 3 {
		t.Errorf("wrong count, expected 3 and got %d", size)
	}
	set.Add("item4") //should not add it, already there
	if size := set.Size(); size != 4 {
		t.Errorf("wrong count, expected 4 and got %d", size)
	}
}

func TestItemSet_AddSlice(t *testing.T) {
	set := populateSet(3, 0)
	if size := set.Size(); size != 3 {
		t.Errorf("wrong count, expected 3 and got %d", size)
	}

	items := []Item{"item2", "item3", "item4"}
	set.AddSlice(items)

	if set.Size() != 5 {
		t.Errorf("wrong count, expected 5 and got %d", set.Size())
	}

	if !set.Has("item3") || !set.Has("item4") {
		t.Error("failed")
	}
}

func TestClear(t *testing.T) {
	set := populateSet(3, 0)
	set.Clear()
	if size := set.Size(); size != 0 {
		t.Errorf("wrong count, expected 0 and got %d", size)
	}
}

func TestDelete(t *testing.T) {
	set := populateSet(3, 0)
	set.Delete("item2")
	if size := set.Size(); size != 2 {
		t.Errorf("wrong count, expected 2 and got %d", size)
	}
}

func TestHas(t *testing.T) {
	set := populateSet(3, 0)
	has := set.Has("item2")
	if !has {
		t.Errorf("expected item2 to be there")
	}
	set.Delete("item2")
	has = set.Has("item2")
	if has {
		t.Errorf("expected item2 to be removed")
	}
	set.Delete("item1")
	has = set.Has("item1")
	if has {
		t.Errorf("expected item1 to be removed")
	}
}

func TestItems(t *testing.T) {
	set := populateSet(3, 0)
	items := set.Items()
	if len(items) != 3 {
		t.Errorf("wrong count, expected 3 and got %d", len(items))
	}
	set = populateSet(520, 0)
	items = set.Items()
	if len(items) != 520 {
		t.Errorf("wrong count, expected 520 and got %d", len(items))
	}
}

func TestSize(t *testing.T) {
	set := populateSet(3, 0)
	items := set.Items()
	if len(items) != set.Size() {
		t.Errorf("wrong count, expected %d and got %d", set.Size(), len(items))
	}
	set = populateSet(0, 0)
	items = set.Items()
	if len(items) != set.Size() {
		t.Errorf("wrong count, expected %d and got %d", set.Size(), len(items))
	}
	set = populateSet(10000, 0)
	items = set.Items()
	if len(items) != set.Size() {
		t.Errorf("wrong count, expected %d and got %d", set.Size(), len(items))
	}
}

func BenchmarkItemSet_Add(b *testing.B) {
	set := populateSet(3, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := rand.Intn(10) + 1
		n := rand.Intn(m)
		addSet(set, m, n)
	}
}
