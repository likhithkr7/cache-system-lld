package cache

import (
	"cache-system-lld/internal/models"
	policy "cache-system-lld/internal/policies"
	"errors"
	"sync"
)

type Cache[K comparable, V any] struct {
	mutex    sync.RWMutex
	capacity int
	policy   policy.CachePolicy[K, V]
	storage  map[K]*models.Entry[K, V]
}

type PolicyType string

const (
	FifoPolicy    PolicyType = "fifo"
	LruPolicyType PolicyType = "lru"
	LfuPolicyType PolicyType = "lfu"
)

func NewCache[K comparable, V any](capacity int, policyType PolicyType) (*Cache[K, V], error) {
	var p policy.CachePolicy[K, V]
	switch policyType {
	case FifoPolicy:
		p = policy.NewFifoPolicy[K, V]()
	case LruPolicyType:
		p = policy.NewLruPolicy[K, V]()
	case LfuPolicyType:
		p = policy.NewLfuPolicy[K, V]()
	default:
		return nil, ErrInvalidPolicy
	}
	if capacity < 1 {
		return nil, ErrInvalidCapacity
	}
	return &Cache[K, V]{
		mutex:    sync.RWMutex{},
		capacity: capacity,
		policy:   p,
		storage:  make(map[K]*models.Entry[K, V]),
	}, nil
}

var (
	ErrInvalidCapacity = errors.New("cache capacity must be greater than 0")
	ErrInvalidPolicy   = errors.New("invalid or unsupported eviction policy type")
	ErrEvictionFailure = errors.New("eviction failure")
)

func (c *Cache[K, V]) Put(key K, value V) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, ok := c.storage[key]; ok && entry != nil {
		entry.Value = value
		c.policy.OnAccess(entry)
		return nil
	}

	if c.policy.Len() >= c.capacity {
		evictedKey, isEvicted := c.policy.OnEvict()
		if isEvicted {
			delete(c.storage, evictedKey)
		} else {
			return ErrEvictionFailure
		}
	}

	newEntry := &models.Entry[K, V]{
		Key:   key,
		Value: value,
	}

	c.policy.OnAdd(newEntry)
	c.storage[key] = newEntry
	return nil
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	entry, ok := c.storage[key]
	c.mutex.RUnlock()

	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if currentEntry, ok := c.storage[key]; !ok || currentEntry != entry {
		var zeroValue V
		return zeroValue, false
	}

	c.policy.OnAccess(entry)
	return entry.Value, true
}

func (c *Cache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, ok := c.storage[key]; ok {
		c.policy.OnRemove(entry)
		delete(c.storage, key)
	}
}

func (c *Cache[K, V]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.policy.Clear()
	for key := range c.storage {
		delete(c.storage, key)
	}
}

func (c *Cache[K, V]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.policy.Len()
}
