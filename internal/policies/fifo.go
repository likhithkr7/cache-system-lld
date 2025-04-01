package policies

import (
	"cache-system-lld/internal/models"
	"container/list"
)

type FifoPolicy[K comparable, V any] struct {
	ll      *list.List
	nodeMap map[K]*list.Element
}

func NewFifoPolicy[K comparable, V any]() *FifoPolicy[K, V] {
	return &FifoPolicy[K, V]{
		ll:      list.New(),
		nodeMap: make(map[K]*list.Element),
	}
}

func (f *FifoPolicy[K, V]) OnAdd(entry *models.Entry[K, V]) {
	element := f.ll.PushFront(entry)
	f.nodeMap[entry.Key] = element
}

func (f *FifoPolicy[K, V]) OnAccess(entry *models.Entry[K, V]) {
	// FIFO: Access does not affect eviction order. Do nothing.
}

func (f *FifoPolicy[K, V]) OnEvict() (K, bool) {
	if f.ll.Len() == 0 {
		var zeroValue K
		return zeroValue, false
	}

	element := f.ll.Back()
	if element != nil {
		entry := element.Value.(*models.Entry[K, V])
		f.OnRemove(entry)
		return entry.Key, true
	}

	var zeroValue K
	return zeroValue, false
}

func (f *FifoPolicy[K, V]) OnRemove(entry *models.Entry[K, V]) {
	if element, ok := f.nodeMap[entry.Key]; ok {
		delete(f.nodeMap, entry.Key)
		f.ll.Remove(element)
	}
}

func (f *FifoPolicy[K, V]) Len() int {
	return f.ll.Len()
}

func (f *FifoPolicy[K, V]) Clear() {
	f.ll.Init()
	for key := range f.nodeMap {
		delete(f.nodeMap, key)
	}
}
