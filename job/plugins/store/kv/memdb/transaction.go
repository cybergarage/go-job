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

package memdb

import (
	"context"
	"errors"
	"time"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/hashicorp/go-memdb"
)

// transaction represents a Memdb transaction instance.
type transaction struct {
	*memdb.Txn
}

func newTransaction(txn *memdb.Txn) *transaction {
	return &transaction{
		Txn: txn,
	}
}

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (txn *transaction) Set(ctx context.Context, obj kv.Object) error {
	return txn.Txn.Insert(
		tableName,
		&Object{
			Key:   obj.Key().Bytes(),
			Value: obj.Bytes(),
		},
	)
}

// Get returns a key-value object of the specified key.
func (txn *transaction) Get(ctx context.Context, key kv.Key) (kv.Object, error) {
	it, err := txn.Txn.Get(tableName, idName, key.Bytes())
	if err != nil {
		return nil, err
	}
	rs := newResultSetWith(it)
	if !rs.Next() {
		return nil, kv.NewErrObjectNotExist(key)
	}
	return rs.Object()
}

// GetRange returns a result set of the specified key.
func (txn *transaction) GetRange(ctx context.Context, key kv.Key, opts ...kv.Option) (kv.ResultSet, error) {
	it, err := txn.Txn.Get(tableName, idName+prefix, key.Bytes())
	if err != nil {
		return nil, err
	}
	return newResultSetWith(it), nil
}

// Remove removes the specified key-value object.
func (txn *transaction) Remove(ctx context.Context, key kv.Key) error {
	obj, err := txn.Get(ctx, key)
	if err != nil {
		return err
	}
	err = txn.Txn.Delete(
		tableName,
		&Object{
			Key:   key.Bytes(),
			Value: obj.Bytes(),
		},
	)
	if err != nil {
		if errors.Is(err, memdb.ErrNotFound) {
			return kv.NewErrObjectNotExist(key)
		}
		return err
	}
	return nil
}

// RemoveRange removes the specified key-value object.
func (txn *transaction) RemoveRange(ctx context.Context, key kv.Key) error {
	_, err := txn.Txn.DeleteAll(tableName, idName+prefix, key.Bytes())
	if err != nil {
		if errors.Is(err, memdb.ErrNotFound) {
			return kv.NewErrObjectNotExist(key)
		}
		return err
	}
	return nil
}

// Commit commits this transaction.
func (txn *transaction) Commit(ctx context.Context) error {
	txn.Txn.Commit()
	return nil
}

// Cancel cancels this transaction.
func (txn *transaction) Cancel(ctx context.Context) error {
	txn.Txn.Abort()
	return nil
}

// SetTimeout sets the timeout of this transaction.
func (txn *transaction) SetTimeout(t time.Duration) error {
	return nil
}
