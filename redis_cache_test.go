// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ Cache = &RedisCache{} // Test for interface compliance

func TestNewRedisCache(t *testing.T) {
	c, err := NewRedisCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)
}

func TestRedisCachePutItem(t *testing.T) {
	c, err := NewRedisCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"))
	assert.NoError(t, err)

	value, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value"), value)
}

func TestRedisCachePutIfAbsent(t *testing.T) {
	c, err := NewRedisCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Evict("key")
	assert.NoError(t, err)

	err = c.PutIfAbsent("key", []byte("value1"))
	assert.NoError(t, err)

	err = c.PutIfAbsent("key", []byte("value2"))
	assert.NoError(t, err)

	value, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value1"), value)
}

func TestRedisCacheExpiredItem(t *testing.T) {
	config := RedisCacheConfig{
		Addr:   "127.0.0.1:6379",
		TTL:    time.Millisecond * 500,
		Prefix: "RedisCache",
	}

	c, err := NewRedisCacheWithConfig(config)
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"))
	assert.NoError(t, err)

	value1, err := c.Get("key")
	assert.NoError(t, err)
	assert.EqualValues(t, []byte("value"), value1)

	time.Sleep(time.Millisecond * 750)

	value2, err := c.Get("key")
	assert.NoError(t, err)
	assert.Nil(t, value2)
}

func TestRedisCacheEvict(t *testing.T) {
	c, err := NewRedisCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"))
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

func TestRedisCacheClear(t *testing.T) {
	c, err := NewRedisCache()
	assert.NotNil(t, c)
	assert.NoError(t, err)

	err = c.Put("key", []byte("value"))
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
