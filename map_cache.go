// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package cache

import (
	"sync"
	"time"
)

type MapCacheConfig struct {
	TTL time.Duration
}

func newMapCacheConfig() MapCacheConfig {
	return MapCacheConfig{
		TTL: time.Minute * 5,
	}
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
	config MapCacheConfig
	values sync.Map
}

func NewMapCache() (*MapCache, error) {
	return NewMapCacheWithConfig(newMapCacheConfig())
}

func NewMapCacheWithConfig(config MapCacheConfig) (*MapCache, error) {
	return &MapCache{
		config: config,
		values: sync.Map{},
	}, nil
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

func (c *MapCache) Put(key string, value []byte) error {
	c.values.Store(key, newMapCacheEntry(value, c.config.TTL))
	return nil
}

func (c *MapCache) PutIfAbsent(key string, value []byte) error {
	c.values.LoadOrStore(key, newMapCacheEntry(value, c.config.TTL))
	return nil
}

func (c *MapCache) Close() error {
	// Stop the collector go routine
	return nil
}
