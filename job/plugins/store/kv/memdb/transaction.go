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

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/hashicorp/go-memdb"
)

// Set stores a key-value object. If the key already holds some value, it is overwritten.
func (db *Database) Set(ctx context.Context, obj kv.Object) error {
	if db == nil {
		return kv.ErrNotReady
	}
	txn := db.Txn(true)
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
	if db == nil {
		return nil, kv.ErrNotReady
	}
	txn := db.Txn(false)
	obj, err := db.get(ctx, txn, key)
	if err != nil {
		txn.Abort()
		return nil, err
	}
	txn.Commit()
	return obj, nil
}

// Scan returns a result set of all key-value objects whose keys have the specified prefix.
func (db *Database) Scan(ctx context.Context, key kv.Key, opts ...kv.Option) (kv.ResultSet, error) {
	if db == nil {
		return nil, kv.ErrNotReady
	}
	txn := db.Txn(false)
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

// Remove removes the specified key-value object.
func (db *Database) Remove(ctx context.Context, obj kv.Object) error {
	if db == nil {
		return kv.ErrNotReady
	}
	key := obj.Key()
	txn := db.Txn(true)
	getObj, err := db.get(ctx, txn, key)
	if err != nil {
		txn.Abort()
		return err
	}
	if !obj.Equal(getObj) {
		return kv.NewErrObjectNotExist(key)
	}
	err = db.remove(ctx, txn, getObj)
	if err != nil {
		txn.Abort()
		return err
	}
	txn.Commit()
	return nil
}

// Delete deletes all key-value objects whose keys have the specified prefix.
func (db *Database) Delete(ctx context.Context, key kv.Key) error {
	if db == nil {
		return kv.ErrNotReady
	}
	txn := db.Txn(true)
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

// Dump returns all key-value objects in the store.
func (db *Database) Dump(ctx context.Context) ([]kv.Object, error) {
	if db == nil {
		return nil, kv.ErrNotReady
	}
	txn := db.Txn(false)
	it, err := txn.Get(tableName, idName)
	if err != nil {
		txn.Abort()
		return nil, err
	}
	rs := newResultSetWith(it)
	var objs []kv.Object
	for rs.Next() {
		obj, err := rs.Object()
		if err != nil {
			txn.Abort()
			return nil, err
		}
		objs = append(objs, obj)
	}
	txn.Commit()
	return objs, nil
}
