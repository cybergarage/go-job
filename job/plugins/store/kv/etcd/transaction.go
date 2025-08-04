// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"context"
	"time"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	v3 "go.etcd.io/etcd/client/v3"
)

// transaction represents a Memdb transaction instance.
type transaction struct {
	write bool
	*v3.Client
	v3.Txn
}

func newTransaction(client *v3.Client, write bool) *transaction {
	return &transaction{
		Client: client,
		write:  write,
		Txn:    client.Txn(context.Background()),
	}
}

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (txn *transaction) Set(ctx context.Context, obj kv.Object) error {
	return nil
}

// Get returns a key-value object of the specified key.
func (txn *transaction) Get(ctx context.Context, key kv.Key) (kv.Object, error) {
	return nil, nil
}

// GetRange returns a result set of the specified key.
func (txn *transaction) GetRange(ctx context.Context, key kv.Key, opts ...kv.Option) (kv.ResultSet, error) {
	return nil, nil
}

// Remove removes the specified key-value object.
func (txn *transaction) Remove(ctx context.Context, key kv.Key) error {
	return nil
}

// RemoveRange removes the specified key-value object.
func (txn *transaction) RemoveRange(ctx context.Context, key kv.Key) error {
	return nil
}

// Commit commits this transaction.
func (txn *transaction) Commit(ctx context.Context) error {
	return nil
}

// Cancel cancels this transaction.
func (txn *transaction) Cancel(ctx context.Context) error {
	return nil
}

// SetTimeout sets the timeout of this transaction.
func (txn *transaction) SetTimeout(t time.Duration) error {
	return nil
}
