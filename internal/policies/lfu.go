package policies

import (
	"cache-system-lld/internal/models"
	"container/list"
)

type freqEntry[K comparable, V any] struct {
	entry *models.Entry[K, V]
	freq  int
}

type LfuPolicy[K comparable, V any] struct {
	ll      *list.List
	nodeMap map[K]*list.Element
	maxFreq int
}

func NewLfuPolicy[K comparable, V any]() *LfuPolicy[K, V] {
	return &LfuPolicy[K, V]{
		ll:      list.New(),
		nodeMap: make(map[K]*list.Element),
		maxFreq: 0,
	}
}

func (l *LfuPolicy[K, V]) OnAdd(entry *models.Entry[K, V]) {
	fEntry := &freqEntry[K, V]{
		entry: entry,
		freq:  1,
	}
	l.maxFreq = max(l.maxFreq, 1)
	element := l.ll.PushFront(fEntry)
	l.nodeMap[entry.Key] = element
}

func (l *LfuPolicy[K, V]) OnAccess(entry *models.Entry[K, V]) {
	if element, ok := l.nodeMap[entry.Key]; ok && element != nil {
		fEntry := element.Value.(*freqEntry[K, V])
		fEntry.freq++
		l.maxFreq = max(l.maxFreq, fEntry.freq)
	}
}

func (l *LfuPolicy[K, V]) OnEvict() (K, bool) {
	if l.ll.Len() == 0 {
		var zeroValue K
		return zeroValue, false
	}

	var evictionElement *list.Element
	minFreq := l.maxFreq
	for _, element := range l.nodeMap {
		if element.Value.(*freqEntry[K, V]).freq < minFreq {
			minFreq = element.Value.(*freqEntry[K, V]).freq
			evictionElement = element
		}
	}
	if evictionElement != nil {
		entry := evictionElement.Value.(*freqEntry[K, V]).entry
		l.OnRemove(entry)
		return entry.Key, true
	}

	var zeroValue K
	return zeroValue, false
}

func (l *LfuPolicy[K, V]) OnRemove(entry *models.Entry[K, V]) {
	if element, ok := l.nodeMap[entry.Key]; ok {
		delete(l.nodeMap, entry.Key)
		l.ll.Remove(element)
	}
}

func (l *LfuPolicy[K, V]) Len() int {
	return l.ll.Len()
}

func (l *LfuPolicy[K, V]) Clear() {
	l.ll.Init()
	for key := range l.nodeMap {
		delete(l.nodeMap, key)
	}
}
