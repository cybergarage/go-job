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

// StoreOption is an alias for valkey.ClientOption, used for configuring the Valkey store.
type StoreOption = valkey.ClientOption

// Store represents a Memdb store service instance.
type Store struct {
	kv.Config
	cofig StoreOption
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
	return nil
}

// Stop stops this memdb.
func (store *Store) Stop() error {
	return nil
}

// Transact returns a new transaction instance.
func (store *Store) Transact(ctx context.Context, write bool) (kv.Transaction, error) {
	return newTransaction(store.Client, write), nil
}
