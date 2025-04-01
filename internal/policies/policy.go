package policies

import "cache-system-lld/internal/models"

type CachePolicy[K comparable, V any] interface {
	OnAccess(entry *models.Entry[K, V])
	OnAdd(entry *models.Entry[K, V])
	OnRemove(entry *models.Entry[K, V])
	OnEvict() (key K, evicted bool)
	Len() int
	Clear()
}
