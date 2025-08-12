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

package etcd

import (
	"context"
	"errors"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kvutil"
	v3 "go.etcd.io/etcd/client/v3"
)

// Store represents a etcd store service instance.
type Store struct {
	kv.Config
	*v3.Client

	opt StoreOption
}

// NewStore returns a new etcd store instance.
func NewStore(option StoreOption) kv.Store {
	return &Store{
		Config: kv.NewConfig(
			kv.WithUniqueKeys(true),
		),
		opt:    option,
		Client: nil,
	}
}

// Name returns the name of this etcd store.
func (store *Store) Name() string {
	return "etcd"
}

// Start starts this etcd.
func (store *Store) Start() error {
	var err error
	store.Client, err = v3.New(store.opt)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), store.opt.DialTimeout)
	defer cancel()
	_, err = store.Client.Status(ctx, store.opt.Endpoints[0])
	if err != nil {
		return errors.Join(err, store.Stop())
	}
	return nil
}

// Clear removes all key-value objects from the store.
func (store *Store) Clear() error {
	if store.Client == nil {
		return nil
	}
	_, err := store.Client.Delete(context.Background(), "", v3.WithPrefix())
	return err
}

// Stop stops this etcd.
func (store *Store) Stop() error {
	if store.Client == nil {
		return nil
	}
	err := store.Client.Close()
	if err != nil {
		return err
	}
	store.Client = nil
	return err
}

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (store *Store) Set(ctx context.Context, obj kv.Object) error {
	if store.Client == nil {
		return kv.ErrNotReady
	}
	_, err := store.Client.Put(ctx, obj.Key().String(), string(obj.Bytes()))
	return err
}

// Get returns a key-value object of the specified key.
func (store *Store) Get(ctx context.Context, key kv.Key) (kv.Object, error) {
	if store.Client == nil {
		return nil, kv.ErrNotReady
	}
	resp, err := store.Client.Get(ctx, key.String())
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, kv.ErrNotExist
	}
	return kv.NewObject(
		kv.Key(resp.Kvs[0].Key),
		resp.Kvs[0].Value), nil
}

// Scan returns a result set of all key-value objects whose keys have the specified prefix.
func (store *Store) Scan(ctx context.Context, key kv.Key, opts ...kv.Option) (kv.ResultSet, error) {
	if store.Client == nil {
		return nil, kv.ErrNotReady
	}
	resp, err := store.Client.Get(ctx, key.String(), v3.WithPrefix())
	if err != nil {
		return nil, err
	}
	objs := []kv.Object{}
	for _, kvs := range resp.Kvs {
		objs = append(objs,
			kv.NewObject(
				kv.Key(kvs.Key),
				kvs.Value,
			))
	}
	return kvutil.NewResultSetWithObjects(objs), nil
}

// Remove removes the specified key-value object.
func (store *Store) Remove(ctx context.Context, obj kv.Object) error {
	if store.Client == nil {
		return kv.ErrNotReady
	}
	_, err := store.Client.Delete(ctx, obj.Key().String())
	return err
}

// Delete deletes all key-value objects whose keys have the specified prefix.
func (store *Store) Delete(ctx context.Context, key kv.Key) error {
	if store.Client == nil {
		return kv.ErrNotReady
	}
	_, err := store.Client.Delete(ctx, key.String(), v3.WithPrefix())
	return err
}

// Dump returns all key-value objects in the store.
func (store *Store) Dump(ctx context.Context) ([]kv.Object, error) {
	if store.Client == nil {
		return nil, kv.ErrNotReady
	}
	resp, err := store.Client.Get(ctx, "", v3.WithPrefix())
	if err != nil {
		return nil, err
	}
	objs := []kv.Object{}
	for _, kvs := range resp.Kvs {
		objs = append(objs,
			kv.NewObject(
				kv.Key(kvs.Key),
				kvs.Value,
			))
	}
	return objs, nil
}
