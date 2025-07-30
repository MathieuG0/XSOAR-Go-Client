package cache

import (
	"context"
	"sync"
	"time"
)

type CacheEntry struct {
	key          string
	data         []byte
	store        *cacheStore
	previous     *CacheEntry
	next         *CacheEntry
	expiresAfter time.Time
	ctx          context.Context
	cancelFunc   context.CancelFunc
}

type cacheStore struct {
	head       *CacheEntry
	tail       *CacheEntry
	locker     sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type Cache struct {
	stores map[string]*cacheStore
	locker sync.RWMutex
}

func (e *CacheEntry) remove() {
	if e.previous == nil {
		if e.store.head == e {
			e.store.head = e.next
		}
	} else {
		e.previous.next = e.next
	}

	if e.next == nil {
		if e.store.tail == e {
			e.store.tail = e.previous
		}
	} else {
		e.next.previous = e.previous
	}
}

func (e *CacheEntry) start() {
	select {
	case <-e.ctx.Done():
	case <-time.Tick(time.Until(e.expiresAfter)):
		e.store.locker.Lock()
		defer e.store.locker.Unlock()
		e.remove()
	}
}

func (e *CacheEntry) Expiration() time.Time {
	return e.expiresAfter
}

func (e *CacheEntry) Data() []byte {
	return e.data
}

func newCacheStore() *cacheStore {
	s := &cacheStore{}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
	return s
}

func (s *cacheStore) removeEntry(key string) {
	cursor := s.head
	for cursor != nil {
		if cursor.key == key {
			cursor.cancelFunc()
			cursor.remove()
		}
	}
}

func (s *cacheStore) addEntry(key string, data []byte, timeout time.Duration) {
	s.locker.Lock()
	defer s.locker.Unlock()

	cursor := s.head
	for cursor != nil {
		if cursor.key != key {
			continue
		}
		cursor.cancelFunc()
		cursor.data = data
		cursor.expiresAfter = time.Now().Add(timeout)
		cursor.ctx, cursor.cancelFunc = context.WithCancel(s.ctx)
		go cursor.start()
		return
	}

	s.removeEntry(key)
	s.tail.next = &CacheEntry{
		key:          key,
		data:         data,
		store:        s,
		previous:     s.tail,
		expiresAfter: time.Now().Add(timeout),
	}
	s.tail.next.ctx, s.tail.next.cancelFunc = context.WithCancel(s.ctx)
	s.tail = s.tail.next
	go s.tail.start()
}

func (s *cacheStore) clear() {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.head = nil
	s.tail = nil
	s.cancelFunc()
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())
}

func (s *cacheStore) get(key string) *CacheEntry {
	s.locker.RLock()
	defer s.locker.RUnlock()
	cursor := s.head
	for cursor != nil {
		if cursor.key == key {
			return cursor
		}
	}
	return nil
}

func NewCache() *Cache {
	return &Cache{
		stores: make(map[string]*cacheStore),
	}
}

func (c *Cache) Add(store string, key string, data []byte, timeout time.Duration) {
	c.locker.Lock()
	defer c.locker.Unlock()

	existing, ok := c.stores[store]
	if ok {
		existing.addEntry(key, data, timeout)
		return
	}

	c.stores[store] = newCacheStore()
	c.stores[store].addEntry(key, data, timeout)
}

func (c *Cache) Get(store string, key string) *CacheEntry {
	c.locker.RLock()
	defer c.locker.RUnlock()

	existing, ok := c.stores[store]
	if !ok {
		return nil
	}

	return existing.get(key)
}

func (c *Cache) Clear(store string) {
	c.locker.Lock()
	defer c.locker.Unlock()

	existing, ok := c.stores[store]
	if ok {
		existing.clear()
	}
}
