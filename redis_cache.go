// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package cache

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisCacheConfig struct {
	Addr   string
	TTL    time.Duration
	Prefix string
}

func newRedisCacheConfig() RedisCacheConfig {
	return RedisCacheConfig{
		Addr:   "127.0.0.1:6379",
		TTL:    time.Minute * 5,
		Prefix: "RedisCache",
	}
}

type RedisCache struct {
	config RedisCacheConfig
	client *redis.Client
}

func NewRedisCache() (*RedisCache, error) {
	return NewRedisCacheWithConfig(newRedisCacheConfig())
}

func NewRedisCacheWithConfig(config RedisCacheConfig) (*RedisCache, error) {
	return &RedisCache{
		config: config,
		client: redis.NewClient(&redis.Options{Addr: config.Addr}),
	}, nil
}

func (c *RedisCache) keyName(key string) string {
	return c.config.Prefix + ":" + key
}

func (c *RedisCache) Clear() error {
	return c.client.Eval("return redis.call('del', unpack(redis.call('keys', ARGV[1])))", nil, c.config.Prefix+":*").Err()
}

func (c *RedisCache) Evict(key string) error {
	return c.client.Del(c.keyName(key)).Err()
}

func (c *RedisCache) Get(key string) ([]byte, error) {
	value, err := c.client.Get(c.keyName(key)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return value, err
}

func (c *RedisCache) Put(key string, value []byte) error {
	return c.client.Set(c.keyName(key), value, c.config.TTL).Err()
}

func (c *RedisCache) PutIfAbsent(key string, value []byte) error {
	return c.client.SetNX(c.keyName(key), value, c.config.TTL).Err()
}

func (c *RedisCache) Close() error {
	return nil
}
