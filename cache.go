// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package cache

import "time"

type Cache interface {
	Clear() error
	Evict(key string) error
	Get(key string) ([]byte, error)
	Put(key string, value []byte, ttl time.Duration) error
	PutIfAbsent(key string, value []byte, ttl time.Duration) error
	Close() error
}
