// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kvutil"
	redis "github.com/redis/go-redis/v9"
)

// Store represents a Memdb store service instance.
type Store struct {
	kv.Config
	*redis.Client

	opt StoreOption
}

// NewStore returns a new memdb store instance.
func NewStore(option StoreOption) kv.Store {
	return &Store{
		Config: kv.NewConfig(
			kv.WithUniqueKeys(false), // Use list commands
		),
		Client: nil,
		opt:    option,
	}
}

// Name returns the name of this memdb store.
func (store *Store) Name() string {
	return "redis"
}

// Start starts this memdb.
func (store *Store) Start() error {
	store.Client = redis.NewClient(&store.opt)
	stat := store.Client.Ping(context.Background())
	if stat.Err() != nil {
		return stat.Err()
	}
	return nil
}

// Stop stops this memdb.
func (store *Store) Stop() error {
	return nil
}

// Clear removes all key-value objects from the store.
func (store *Store) Clear() error {
	if store.Client == nil {
		return kv.ErrNotReady
	}
	stat := store.Client.FlushAll(context.Background())
	if stat.Err() != nil {
		return stat.Err()
	}
	return nil
}

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (store *Store) Set(ctx context.Context, obj kv.Object) error {
	if store.Client == nil {
		return kv.ErrNotReady
	}
	listKey := obj.Key().String()
	stat := store.Client.RPush(ctx, listKey, string(obj.Bytes()))
	if stat.Err() != nil {
		return stat.Err()
	}
	return nil
}

// Range returns a list of values for the specified key.
func (store *Store) Range(ctx context.Context, key kv.Key, limit int64) ([]string, error) {
	if store.Client == nil {
		return nil, kv.ErrNotReady
	}
	listKey := key.String()
	stat := store.Client.LRange(ctx, listKey, 0, limit)
	if stat.Err() != nil {
		return nil, stat.Err()
	}
	return stat.Val(), nil
}

// Get returns a key-value object of the specified key.
func (store *Store) Get(ctx context.Context, key kv.Key) (kv.Object, error) {
	elems, err := store.Range(ctx, key, 0)
	if err != nil {
		return nil, err
	}
	if len(elems) == 0 {
		return nil, kv.ErrNotExist
	}
	return kv.NewObject(key, []byte(elems[0])), nil
}

// Scan returns a result set of all key-value objects whose keys have the specified prefix.
func (store *Store) Scan(ctx context.Context, key kv.Key, opts ...kv.Option) (kv.ResultSet, error) {
	elems, err := store.Range(ctx, key, -1)
	if err != nil {
		return nil, err
	}
	objs := make([]kv.Object, len(elems))
	for i, elem := range elems {
		objs[i] = kv.NewObject(key, []byte(elem))
	}
	return kvutil.NewResultSetWithObjects(objs), nil
}

// Remove removes the specified key-value object.
func (store *Store) Remove(ctx context.Context, obj kv.Object) error {
	if store.Client == nil {
		return kv.ErrNotReady
	}
	key := obj.Key()
	listKey := key.String()
	listValue := string(obj.Bytes())
	stat := store.Client.LRem(ctx, listKey, 0, listValue)
	if stat.Err() != nil {
		return stat.Err()
	}
	cnt := stat.Val()
	if cnt < 1 {
		return kv.NewErrKeyObjectNotExist(key)
	}
	return nil
}

// Delete deletes all key-value objects whose keys have the specified prefix.
func (store *Store) Delete(ctx context.Context, key kv.Key) error {
	if store.Client == nil {
		return kv.ErrNotReady
	}
	listKey := key.String()
	stat := store.Client.Del(ctx, listKey)
	if stat.Err() != nil {
		return stat.Err()
	}
	return nil
}

// Dump returns all key-value objects in the store.
func (store *Store) Dump(ctx context.Context) ([]kv.Object, error) {
	if store.Client == nil {
		return nil, kv.ErrNotReady
	}

	keys, err := store.Client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	objs := []kv.Object{}
	for _, key := range keys {
		obj, err := store.Get(ctx, kv.Key(key))
		if err != nil {
			return nil, err
		}
		objs = append(objs, obj)
	}

	return objs, nil
}
