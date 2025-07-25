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

// Key represents a key in the key-value store.
type Key string

// Object represents a key-value object.
type Object interface {
	// Key returns a key of the object.
	Key() Key
	// Value returns a value of the object.
	Value() []byte
}

// Option represents a option.
type Option = any

// ResultSet represents a result set which includes query execution results.
type ResultSet interface {
	// Next moves the cursor forward next object from its current position.
	Next() bool
	// Object returns an object in the current cursor.
	Object() (Object, error)
}

// Transaction represents a transaction interface.
type Transaction interface {
	// Set stores a key-value object. If the key already holds some value, it is overwritten.
	Set(obj Object) error
	// Get returns a key-value object of the specified key.
	Get(key Key) (Object, error)
	// GetRange returns a result set of the specified key.
	GetRange(key Key, opts ...Option) (ResultSet, error)
	// Remove removes the specified key-value object.
	Remove(key Key) error
	// RemoveRange removes the specified key-value objects.
	RemoveRange(key Key) error
	// Commit commits this transaction.
	Commit() error
	// Cancel cancels this transaction.
	Cancel() error
}

// Store represents a store interface.
type KvStore interface {
	// Transact begin a new transaction.
	Transact(write bool) (Transaction, error)
}
