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
	"github.com/valkey-io/valkey-go"
)

// Store represents a Memdb store service instance.
type Store struct {
	kv.Config
	valkey.Client
}

// NewStore returns a new memdb store instance.
func NewStore(option StoreOption) (kv.Store, error) {
	client, err := valkey.NewClient(option)
	if err != nil {
		return nil, err
	}
	return &Store{
		Config: kv.NewConfig(
			kv.WithUniqueKeys(false), // Use list commands
		),
		Client: client,
	}, nil
}

// Name returns the name of this memdb store.
func (store *Store) Name() string {
	return "valkey"
}

// Start starts this memdb.
func (store *Store) Start() error {
	cmd := store.B().Hello()
	return store.Do(context.Background(), cmd.Build()).Error()
}

// Stop stops this memdb.
func (store *Store) Stop() error {
	return nil
}

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (store *Store) Set(ctx context.Context, obj kv.Object) error {
	cmd := store.B().Set().Key(obj.Key().String()).Value(string(obj.Bytes()))
	return store.Do(ctx, cmd.Build()).Error()
}

// Get returns a key-value object of the specified key.
func (store *Store) Get(ctx context.Context, key kv.Key) (kv.Object, error) {
	return nil, nil
}

// Scan returns a result set of all key-value objects whose keys have the specified prefix.
func (store *Store) Scan(ctx context.Context, key kv.Key, opts ...kv.Option) (kv.ResultSet, error) {
	return nil, nil
}

// Remove removes and returns the key-value object of the specified key.
func (store *Store) Remove(ctx context.Context, key kv.Key) (kv.Object, error) {
	return nil, nil
}

// Delete deletes all key-value objects whose keys have the specified prefix.
func (store *Store) Delete(ctx context.Context, key kv.Key) error {
	return nil
}
