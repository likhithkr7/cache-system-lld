package policies

import (
	"cache-system-lld/internal/models"
	"container/list"
)

type LruPolicy[K comparable, V any] struct {
	ll      *list.List
	nodeMap map[K]*list.Element
}

func NewLruPolicy[K comparable, V any]() *LruPolicy[K, V] {
	return &LruPolicy[K, V]{
		ll:      list.New(),
		nodeMap: make(map[K]*list.Element),
	}
}

func (l *LruPolicy[K, V]) OnAdd(entry *models.Entry[K, V]) {
	element := l.ll.PushFront(entry)
	l.nodeMap[entry.Key] = element
}

func (l *LruPolicy[K, V]) OnAccess(entry *models.Entry[K, V]) {
	if element, ok := l.nodeMap[entry.Key]; ok {
		l.ll.MoveToFront(element)
	}
}

func (l *LruPolicy[K, V]) OnEvict() (K, bool) {
	if l.ll.Len() == 0 {
		var zeroValue K
		return zeroValue, false
	}

	element := l.ll.Back()
	if element != nil {
		entry := element.Value.(*models.Entry[K, V])
		l.OnRemove(entry)
		return entry.Key, true
	}

	var zeroValue K
	return zeroValue, false
}

func (l *LruPolicy[K, V]) OnRemove(entry *models.Entry[K, V]) {
	if element, ok := l.nodeMap[entry.Key]; ok {
		delete(l.nodeMap, entry.Key)
		l.ll.Remove(element)
	}
}

func (l *LruPolicy[K, V]) Len() int {
	return l.ll.Len()
}

func (l *LruPolicy[K, V]) Clear() {
	l.ll.Init()
	for key := range l.nodeMap {
		delete(l.nodeMap, key)
	}
}
