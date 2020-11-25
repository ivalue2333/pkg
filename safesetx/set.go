package safesetx

import (
	"sync"
)

/*
	A Set is a collection of values. You can iterate over those values, add new values,
remove values and clear the set, get the set size, and check if the set contains an item.
A value in the set might only be stored once, duplicates are not possible.
*/

// Item the type of the Set
type Item interface{}

// ItemSet the set of Items
type ItemSet struct {
	items map[Item]bool
	lock  sync.RWMutex
}

func NewItemSet() *ItemSet {
	return &ItemSet{}
}

// Add adds a new element to the Set. Returns a pointer to the Set.
func (s *ItemSet) Add(t Item) *ItemSet {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items == nil {
		s.items = make(map[Item]bool)
	}
	_, ok := s.items[t]
	if !ok {
		s.items[t] = true
	}
	return s
}

func (s *ItemSet) AddSlice(ts []Item) *ItemSet {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.items == nil {
		s.items = make(map[Item]bool)
	}
	for _, t := range ts {
		_, ok := s.items[t]
		if !ok {
			s.items[t] = true
		}
	}
	return s
}

// Clear removes all elements from the Set
func (s *ItemSet) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.items = make(map[Item]bool)
}

// Delete removes the Item from the Set and returns Has(Item)
func (s *ItemSet) Delete(item Item) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	_, ok := s.items[item]
	if ok {
		delete(s.items, item)
	}
	return ok
}

// Has returns true if the Set contains the Item
func (s *ItemSet) Has(item Item) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.items[item]
	return ok
}

// Items returns the Item(s) stored
func (s *ItemSet) Items() []Item {
	s.lock.RLock()
	defer s.lock.RUnlock()
	items := []Item{}
	for i := range s.items {
		items = append(items, i)
	}
	return items
}

// Size returns the size of the set
func (s *ItemSet) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
