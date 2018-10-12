// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ Cache = &MapCache{} // Test for interface compliance

func TestNewMapCache(t *testing.T) {
	c, err := NewMapCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	assert.EqualValues(t, defaultExpirationInterval, c.config.ExpirationInterval)
}

func TestNewMapCacheWithConfigExpirationInterval(t *testing.T) {
	c, err := NewMapCacheWithConfig(MapCacheConfig{
		ExpirationInterval: time.Minute * 1,
	})
	assert.NotNil(t, c)
	assert.NoError(t, err)
	assert.EqualValues(t, time.Minute*1, c.config.ExpirationInterval)
}

func TestMapCachePutItem(t *testing.T) {
	c, err := NewMapCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"), time.Minute)
	assert.NoError(t, err)

	value, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value"), value)
}

func TestMapCachePutIfAbsent(t *testing.T) {
	c, err := NewMapCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.PutIfAbsent("key", []byte("value1"), time.Minute)
	assert.NoError(t, err)

	err = c.PutIfAbsent("key", []byte("value2"), time.Minute)
	assert.NoError(t, err)

	value, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value1"), value)
}

func TestMapCacheExpiredItem(t *testing.T) {
	c, err := NewMapCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"), time.Millisecond*500)
	assert.NoError(t, err)

	value1, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value"), value1)

	time.Sleep(time.Millisecond * 750)

	value2, err := c.Get("key")
	assert.NoError(t, err)
	assert.Nil(t, value2)
}

func TestMapCacheEvict(t *testing.T) {
	c, err := NewMapCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"), time.Minute)
	assert.NoError(t, err)

	value, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value"), value)

	err = c.Evict("key")
	assert.NoError(t, err)

	value2, err := c.Get("key")
	assert.NoError(t, err)
	assert.Nil(t, value2)

}

func TestMapCacheClear(t *testing.T) {
	c, err := NewMapCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"), time.Minute)
	assert.NoError(t, err)

	value, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value"), value)

	err = c.Clear()
	assert.NoError(t, err)

	value2, err := c.Get("key")
	assert.NoError(t, err)
	assert.Nil(t, value2)
}

func TestExpireEntriesTask(t *testing.T) {
	c, err := NewMapCacheWithConfig(MapCacheConfig{
		ExpirationInterval: time.Millisecond * 250,
	})
	assert.NotNil(t, c)
	assert.NoError(t, err)

	defer c.Close()

	err = c.Put("key", []byte("value"), time.Second)
	assert.NoError(t, err)

	value, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value"), value)

	time.Sleep(time.Millisecond * 2500)

	value2, err := c.Get("key")
	assert.NoError(t, err)
	assert.Nil(t, value2)

	// Directly access the map to make sure the entry has been removed
	val, ok := c.values.Load("key")
	assert.Nil(t, val)
	assert.False(t, ok)
}
