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

package valkey

import (
	"context"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kvutil"
	"github.com/valkey-io/valkey-go"
)

// Store represents a Memdb store service instance.
type Store struct {
	kv.Config
	valkey.Client

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
	return "valkey"
}

// Start starts this memdb.
func (store *Store) Start() error {
	var err error
	store.Client, err = valkey.NewClient(store.opt)
	if err != nil {
		return err
	}
	cmd := store.B().Hello()
	return store.Do(context.Background(), cmd.Build()).Error()
}

// Stop stops this memdb.
func (store *Store) Stop() error {
	return nil
}

// Clear removes all key-value objects from the store.
func (store *Store) Clear() error {
	cmdList := store.B().Flushall()
	err := store.Do(context.Background(), cmdList.Build()).Error()
	if err != nil {
		return err
	}
	return nil
}

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (store *Store) Set(ctx context.Context, obj kv.Object) error {
	listKey := obj.Key().String()
	cmdList := store.B().Rpush().Key(listKey).Element(string(obj.Bytes()))
	err := store.Do(ctx, cmdList.Build()).Error()
	if err != nil {
		return err
	}
	return nil
}

func (store *Store) Range(ctx context.Context, key kv.Key, limit int64) ([]string, error) {
	listKey := key.String()
	cmd := store.B().Lrange().Key(listKey).Start(0).Stop(limit)
	resp := store.Do(ctx, cmd.Build())
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	return resp.AsStrSlice()
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

// Remove removes and returns the key-value object of the specified key.
func (store *Store) Remove(ctx context.Context, key kv.Key) (kv.Object, error) {
	listKey := key.String()
	cmd := store.B().Lpop().Key(listKey)
	resp := store.Do(ctx, cmd.Build())
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	val, err := resp.AsBytes()
	if err != nil {
		return nil, err
	}
	return kv.NewObject(key, []byte(val)), nil
}

// Delete deletes all key-value objects whose keys have the specified prefix.
func (store *Store) Delete(ctx context.Context, key kv.Key) error {
	listKey := key.String()
	cmd := store.B().Del().Key(listKey)
	err := store.Do(ctx, cmd.Build()).Error()
	if err != nil {
		return err
	}
	return nil
}
