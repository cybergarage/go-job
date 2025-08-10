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

import "context"

// Option represents a option.
type Option = any

// Store represents a key-value store interface.
type Store interface {
	// Config defines the store configuration.
	Config
	// Name returns the name of the store.
	Name() string
	// Set stores a key-value object. If the key already holds some value, it is overwritten.
	Set(ctx context.Context, obj Object) error
	// Get returns a key-value object of the specified key.
	Get(ctx context.Context, key Key) (Object, error)
	// Scan returns a result set of all key-value objects whose keys have the specified prefix.
	Scan(ctx context.Context, key Key, opts ...Option) (ResultSet, error)
	// Remove removes the specified key-value object.
	Remove(ctx context.Context, obj Object) error
	// Dump returns all key-value objects in the store.
	Dump(ctx context.Context) ([]Object, error)
	// Delete deletes all key-value objects whose keys have the specified prefix.
	Delete(ctx context.Context, key Key) error
	// Start starts the store.
	Start() error
	// Stop stops the store.
	Stop() error
	// Clear removes all key-value objects from the store.
	Clear() error
}
