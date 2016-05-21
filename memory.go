package bytecache

import (
	"log"
	"sync"
	"time"
)

type Item struct {
	Object     []byte
	Expiration int64
}

type SimpleCache struct {
	items      map[string]Item
	mu         sync.RWMutex
	gcstop     chan bool
	gcinterval int
}

func NewSimpleCache(gcinterval int) *SimpleCache {
	c := &SimpleCache{
		items:      make(map[string]Item),
		gcstop:     make(chan bool),
		gcinterval: gcinterval,
	}
	c.StartGarbageCollector()
	return c
}

func (c *SimpleCache) StopGarbageCollector() {
	c.gcstop <- true
}

func (c *SimpleCache) StartGarbageCollector() {
	// GC jede Stunde
	ticker := time.NewTicker(time.Duration(c.gcinterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.deleteExpired()
			case <-c.gcstop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *SimpleCache) deleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.items, k)
			log.Printf("Garbage collected %s\n", k)
		}
	}
}

func (c *SimpleCache) Set(key string, value []byte, ttl int) error {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(time.Duration(ttl) * time.Second).UnixNano()
	} else {
		exp = -1
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Object:     value,
		Expiration: exp,
	}
	return nil
}

func (c *SimpleCache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found {
		return nil, nil
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, nil
		}
	}
	return item.Object, nil
}

func (c *SimpleCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}
