package calsync

import (
	"iter"
)

type OrderMap[K comparable, V any] struct {
	_map  map[K]V
	items []K
}

func NewOrderMap[K comparable, V any]() (self OrderMap[K, V]) {
	self._map = make(map[K]V)
	return
}
func (self *OrderMap[K, V]) Len() int {
	return len(self.items)
}
func (self *OrderMap[K, V]) addKey(key K) {
	if _, exists := self._map[key]; !exists {
		self.items = append(self.items, key)
	}
}
func (self *OrderMap[K, V]) Has(key K) bool {
	_, exists := self._map[key]
	return exists
}

func (self *OrderMap[K, V]) Get(key K) V {
	return self._map[key]
}

func (self *OrderMap[K, V]) GetOrInsert(key K, val V) V {
	if !self.Has(key) {
		self.Set(key, val)
	}
	return self.Get(key)
}

func (self *OrderMap[K, V]) Set(key K, val V) {
	self.addKey(key)
	self._map[key] = val
}

func (self *OrderMap[K, V]) Update(key K, update func(V) V) {
	self.addKey(key)
	self._map[key] = update(self.Get(key))
}

func (self *OrderMap[K, V]) Items() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range self.items {
			if !yield(k, self._map[k]) {
				return
			}
		}
	}
}
