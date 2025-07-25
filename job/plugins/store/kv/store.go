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

package kv

// KvKey represents a key in the key-value store.
type KvKey string

// KvObject represents a key-value object.
type KvObject interface {
	// Key returns a key of the object.
	Key() KvKey
	// Value returns a value of the object.
	Value() []byte
}

// KvOption represents a option.
type KvOption = any

// KvResultSet represents a result set which includes query execution results.
type KvResultSet interface {
	// Next moves the cursor forward next object from its current position.
	Next() bool
	// Object returns an object in the current cursor.
	Object() (KvObject, error)
}

// KvTransaction represents a transaction interface.
type KvTransaction interface {
	// Set stores a key-value object. If the key already holds some value, it is overwritten.
	Set(obj KvObject) error
	// Get returns a key-value object of the specified key.
	Get(key KvKey) (KvObject, error)
	// GetRange returns a result set of the specified key.
	GetRange(key KvKey, opts ...KvOption) (KvResultSet, error)
	// Remove removes the specified key-value object.
	Remove(key KvKey) error
	// RemoveRange removes the specified key-value objects.
	RemoveRange(key KvKey) error
	// Commit commits this transaction.
	Commit() error
	// Cancel cancels this transaction.
	Cancel() error
}

// Store represents a store interface.
type KvStore interface {
	// Transact begin a new transaction.
	Transact(write bool) (KvTransaction, error)
}
