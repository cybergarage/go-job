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

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (db *Database) Set(ctx context.Context, obj kv.Object) error {
	txn := db.MemDB.Txn(true)
	err := txn.Insert(
		tableName,
		&Object{
			Key:   obj.Key().Bytes(),
			Value: obj.Bytes(),
		},
	)
	if err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}

// Get returns a key-value object of the specified key.
func (db *Database) get(ctx context.Context, txn *memdb.Txn, key kv.Key) (kv.Object, error) {
	it, err := txn.Get(tableName, idName, key.Bytes())
	if err != nil {
		return nil, err
	}
	rs := newResultSetWith(it)
	if !rs.Next() {
		return nil, kv.NewErrObjectNotExist(key)
	}
	return rs.Object()
}

// Get returns a key-value object of the specified key.
func (db *Database) Get(ctx context.Context, key kv.Key) (kv.Object, error) {
	txn := db.MemDB.Txn(false)
	obj, err := db.get(ctx, txn, key)
	if err != nil {
		txn.Abort()
		return nil, err
	}
	txn.Commit()
	return obj, nil
}

// GetRange returns a result set of the specified key.
func (db *Database) GetRange(ctx context.Context, key kv.Key, opts ...kv.Option) (kv.ResultSet, error) {
	txn := db.MemDB.Txn(false)
	it, err := txn.Get(tableName, idName+prefix, key.Bytes())
	if err != nil {
		txn.Abort()
		return nil, err
	}
	return newResultSetWith(it), nil
}

func (db *Database) remove(ctx context.Context, txn *memdb.Txn, kvObj kv.Object) error {
	return txn.Delete(
		tableName,
		&Object{
			Key:   kvObj.Key().Bytes(),
			Value: kvObj.Bytes(),
		},
	)
}

// Remove removes and returns the key-value object of the specified key.
func (db *Database) Remove(ctx context.Context, key kv.Key) (kv.Object, error) {
	txn := db.MemDB.Txn(true)
	obj, err := db.get(ctx, txn, key)
	if err != nil {
		txn.Abort()
		return nil, err
	}
	err = db.remove(ctx, txn, obj)
	if err != nil {
		txn.Abort()
		return nil, err
	}
	txn.Commit()
	return obj, nil
}

// RemoveRange removes the specified key-value object.
func (db *Database) RemoveRange(ctx context.Context, key kv.Key) error {
	txn := db.MemDB.Txn(true)
	_, err := txn.DeleteAll(tableName, idName+prefix, key.Bytes())
	if err != nil {
		txn.Abort()
		if errors.Is(err, memdb.ErrNotFound) {
			return kv.NewErrObjectNotExist(key)
		}
		return err
	}
	txn.Commit()
	return nil
}

// SetTimeout sets the timeout of this transaction.
func (db *Database) SetTimeout(t time.Duration) error {
	return nil
}
