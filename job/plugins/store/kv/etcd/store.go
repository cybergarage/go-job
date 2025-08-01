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

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	v3 "go.etcd.io/etcd/client/v3"
)

// StoreOption is an alias for etcd.ClientOption, used for configuring the etcd store.
type StoreOption = v3.Config

// Store represents a etcd store service instance.
type Store struct {
	kv.Config
	*v3.Client
}

// NewStore returns a new etcd store instance.
func NewStore(option StoreOption) (kv.Store, error) {
	client, err := v3.New(option)
	if err != nil {
		return nil, err
	}
	return &Store{
		Config: kv.NewConfig(
			kv.WithUniqueKeys(true),
		),
		Client: client,
	}, nil
}

// Name returns the name of this etcd store.
func (store *Store) Name() string {
	return "etcd"
}

// Start starts this etcd.
func (store *Store) Start() error {
	return nil
}

// Stop stops this etcd.
func (store *Store) Stop() error {
	if store.Client == nil {
		return nil
	}
	err := store.Client.Close()
	if err == nil {
		store.Client = nil
		return nil
	}
	return err
}

// Transact returns a new transaction instance.
func (store *Store) Transact(ctx context.Context, write bool) (kv.Transaction, error) {
	return newTransaction(store.Client, write), nil
}
