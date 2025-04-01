package main

import (
	"cache-system-lld/internal/cache"
	"fmt"
)

func main() {
	fmt.Println("*** FIFO cache ***")
	fifoCache, err := cache.NewCache[string, int](3, cache.FifoPolicy)
	if err != nil {
		return
	}
	err = fifoCache.Put("X", 10)
	if err != nil {
		return
	}
	fmt.Println("Inserted X")

	err = fifoCache.Put("Y", 20)
	if err != nil {
		return
	}
	fmt.Println("Inserted Y")

	err = fifoCache.Put("Z", 30)
	if err != nil {
		return
	}
	fmt.Println("Inserted Z")

	val, _ := fifoCache.Get("Y")
	fmt.Println("Y found :", val)

	fifoCache.Delete("X")
	fmt.Println("Deleted X")

	err = fifoCache.Put("W", 40)
	if err != nil {
		return
	}
	fmt.Println("Inserted W")

	err = fifoCache.Put("Z", 50)
	if err != nil {
		return
	}
	fmt.Println("Updated Z")

	err = fifoCache.Put("V", 60)
	if err != nil {
		return
	}
	fmt.Println("Inserted V")

	if _, ok := fifoCache.Get("Y"); !ok {
		fmt.Println("Y was evicted (FIFO)")
	}
	if v, ok := fifoCache.Get("Z"); ok {
		fmt.Println("Z found:", v)
	}

	fmt.Println("\n*** LRU cache ***")
	lruCache, err := cache.NewCache[string, int](3, cache.LruPolicyType)
	if err != nil {
		return
	}
	err = lruCache.Put("X", 10)
	if err != nil {
		return
	}
	fmt.Println("Inserted X")

	err = lruCache.Put("Y", 20)
	if err != nil {
		return
	}
	fmt.Println("Inserted Y")

	err = lruCache.Put("Z", 30)
	if err != nil {
		return
	}
	fmt.Println("Inserted Z")

	val, _ = lruCache.Get("Y")
	fmt.Println("Y found: ", val)

	err = lruCache.Put("W", 40)
	if err != nil {
		return
	}
	fmt.Println("Inserted W")

	if _, ok := lruCache.Get("X"); !ok {
		fmt.Println("X was evicted (LRU)")
	}
	if v, ok := lruCache.Get("W"); ok {
		fmt.Println("W found:", v)
	}

	err = lruCache.Put("S", 50)
	if err != nil {
		return
	}
	fmt.Println("Inserted S")

	if _, ok := lruCache.Get("Z"); !ok {
		fmt.Println("Z was evicted (LRU)")
	}

	fmt.Println("\n*** LFU cache ***")
	lfuCache, err := cache.NewCache[string, int](3, cache.LfuPolicyType)
	if err != nil {
		return
	}
	err = lfuCache.Put("X", 10)
	if err != nil {
		return
	}
	fmt.Println("Inserted X")

	err = lfuCache.Put("Y", 20)
	if err != nil {
		return
	}
	fmt.Println("Inserted Y")
	err = lfuCache.Put("Z", 30)
	if err != nil {
		return
	}
	fmt.Println("Inserted Z")

	val, _ = lfuCache.Get("X")
	fmt.Println("X found: ", val)
	val, _ = lfuCache.Get("Z")
	fmt.Println("Z found: ", val)

	err = lfuCache.Put("W", 40)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println("Inserted W")
	val, _ = lfuCache.Get("W")
	fmt.Println("W found: ", val)

	if _, ok := lfuCache.Get("Y"); !ok {
		fmt.Println("Y was evicted (LFU)")
	}

}
