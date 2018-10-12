// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package cache

import (
	"sync"
	"time"

	"github.com/st3fan/task"
)

const (
	defaultExpirationInterval = time.Minute * 1
)

type MapCacheConfig struct {
	ExpirationInterval time.Duration
}

func newMapCacheConfig() MapCacheConfig {
	return MapCacheConfig{
		ExpirationInterval: defaultExpirationInterval,
	}
}

func updateMapCacheConfig(override MapCacheConfig) MapCacheConfig {
	config := newMapCacheConfig()
	if override.ExpirationInterval != 0 {
		config.ExpirationInterval = override.ExpirationInterval
	}
	return config
}

type mapCacheEntry struct {
	value   []byte
	expires time.Time
}

func newMapCacheEntry(value []byte, ttl time.Duration) *mapCacheEntry {
	return &mapCacheEntry{
		value:   value,
		expires: time.Now().Add(ttl),
	}
}

func (e *mapCacheEntry) expired() bool {
	return e.expires.Before(time.Now())
}

type MapCache struct {
	config         MapCacheConfig
	values         sync.Map
	expirationTask *task.Task
}

func NewMapCache() (*MapCache, error) {
	return NewMapCacheWithConfig(newMapCacheConfig())
}

func NewMapCacheWithConfig(config MapCacheConfig) (*MapCache, error) {
	cache := &MapCache{
		config: updateMapCacheConfig(config),
		values: sync.Map{},
	}
	cache.expirationTask = task.New(cache.expireEntries)
	return cache, nil
}

func (c *MapCache) Clear() error {
	c.values = sync.Map{}
	return nil
}

func (c *MapCache) Evict(key string) error {
	c.values.Delete(key)
	return nil
}

func (c *MapCache) Get(key string) ([]byte, error) {
	if value, ok := c.values.Load(key); ok {
		entry := value.(*mapCacheEntry)
		if !entry.expired() {
			return entry.value, nil
		}
	}
	return nil, nil
}

func (c *MapCache) Put(key string, value []byte, ttl time.Duration) error {
	c.values.Store(key, newMapCacheEntry(value, ttl))
	return nil
}

func (c *MapCache) PutIfAbsent(key string, value []byte, ttl time.Duration) error {
	c.values.LoadOrStore(key, newMapCacheEntry(value, ttl))
	return nil
}

func (c *MapCache) Close() error {
	c.expirationTask.SignalAndWait()
	return nil
}

func (c *MapCache) expireEntries(task *task.Task) {
	ticker := time.NewTicker(c.config.ExpirationInterval)

	for {
		select {
		case <-ticker.C:
			c.values.Range(func(key, value interface{}) bool {
				entry := value.(*mapCacheEntry)
				if entry.expired() {
					c.values.Delete(key)
				}
				return false
			})
		case <-task.HasBeenClosed():
			return
		}
	}
}
