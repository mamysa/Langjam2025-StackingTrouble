package util

import (
	"fmt"
	"iter"
	"slices"
)

type OrderedMap[K comparable, V any] struct {
	entries        map[K]V
	keySetOrdering []K
}

func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return OrderedMap[K, V]{
		entries:        make(map[K]V),
		keySetOrdering: make([]K, 0),
	}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	m.entries[key] = value

	if slices.Contains(m.keySetOrdering, key) {
		return
	}

	m.keySetOrdering = append(m.keySetOrdering, key)
}

func (m *OrderedMap[K, V]) Value(key K) (V, error) {
	v, ok := m.entries[key]
	if !ok {
		return *new(V), fmt.Errorf("Key %+v not present in ordered map", key)
	}

	return v, nil
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	_, ok := m.entries[key]
	return ok
}

func (m *OrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, key := range m.keySetOrdering {
			val := m.entries[key]
			if !yield(key, val) {
				return
			}
		}
	}
}

func (m *OrderedMap[K, V]) String() string {
	return fmt.Sprintf("OrderedMap(entries: %+v, orderedKeySet: %+v)", m.entries, m.keySetOrdering)
}
