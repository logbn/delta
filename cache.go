package mvfifo

import (
	"container/list"
	"iter"
	"sync"
)

var (
	// DefaultMaxSizeBytes sets the default max size at 256 MiB
	DefaultMaxSizeBytes = 256 << 20

	overhead = 16
	itemPool = sync.Pool{
		New: func() any {
			return new(item)
		},
	}
)

type item struct {
	key string
	cur uint64
	val []byte
}

// Cache implements a multi value FIFO cache.
type Cache struct {
	items   *list.List
	maxSize int
	mutex   sync.RWMutex
	size    int
	vals    map[string]*list.List
}

// New returns a new [Cache].
func NewCache(opts ...Option) *Cache {
	c := &Cache{
		maxSize: DefaultMaxSizeBytes,
		vals:    map[string]*list.List{},
		items:   list.New(),
	}
	for _, fn := range opts {
		fn(c)
	}
	return c
}

// Add adds an item to the cache by key, cursor and value.
func (c *Cache) Add(key string, cur uint64, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.size += len(key) + len(val) + overhead
	for c.size > c.maxSize {
		c.evict()
	}
	kl, ok := c.vals[key]
	if !ok {
		kl = list.New()
		c.vals[key] = kl
	}
	i := itemPool.Get().(*item)
	i.key = key
	i.cur = cur
	i.val = val
	kl.PushBack(i)
	c.items.PushBack(i)
}

// Resize changes the maximum size of the cache.
func (c *Cache) Resize(maxBytes int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.maxSize = maxBytes
	for c.size > c.maxSize {
		c.evict()
	}
}

// Size returns the approximate size of the cache in bytes.
func (c *Cache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.size
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.items.Len()
}

// First returns the oldest cursor and value in the cache.
func (c *Cache) First() (cur uint64, val []byte) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if el := c.items.Front(); el != nil {
		cur = el.Value.(*item).cur
		val = el.Value.(*item).val
	}
	return
}

// Last returns the newest cursor and value in the cache.
func (c *Cache) Last() (cur uint64, val []byte) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if el := c.items.Back(); el != nil {
		cur = el.Value.(*item).cur
		val = el.Value.(*item).val
	}
	return
}

// Iter returns an iterator over the cache at a certain key.
func (c *Cache) Iter(key string) iter.Seq2[uint64, []byte] {
	return func(yield func(cur uint64, val []byte) bool) {
		c.mutex.RLock()
		defer c.mutex.RUnlock()
		kl, ok := c.vals[key]
		if !ok {
			return
		}
		var i *item
		for el := kl.Front(); el != nil; el = el.Next() {
			i = el.Value.(*item)
			if !yield(i.cur, i.val) {
				return
			}
		}
	}
}

// IterAfter returns an iterator over the cache at a certain key after a specific cursor.
// Complexity is O(2n). We iterate backward over the list then forward to optimize for access to later values.
func (c *Cache) IterAfter(key string, cur uint64) iter.Seq2[uint64, []byte] {
	return func(yield func(cur uint64, val []byte) bool) {
		c.mutex.RLock()
		defer c.mutex.RUnlock()
		kl, ok := c.vals[key]
		if !ok {
			return
		}
		var i *item
		el := kl.Back()
		for ; el != nil && el.Value.(*item).cur > cur; el = el.Prev() {
		}
		if el == nil {
			el = kl.Front()
		}
		for ; el != nil; el = el.Next() {
			i = el.Value.(*item)
			if i.cur > cur && !yield(i.cur, i.val) {
				return
			}
		}
	}
}

func (c *Cache) evict() {
	el := c.items.Front()
	item := el.Value.(*item)
	defer itemPool.Put(item)
	kl := c.vals[item.key]
	if kl.Len() == 1 {
		delete(c.vals, item.key)
	} else {
		kl.Remove(kl.Front())
	}
	c.items.Remove(el)
	c.size -= len(item.key) + len(item.val) + overhead
}
