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

import (
	"context"
)

// Transaction represents a transaction interface.
type Transaction interface {
	// Set stores a key-value object. If the key already holds some value, it is overwritten.
	Set(ctx context.Context, obj Object) error
	// Get returns a key-value object of the specified key.
	Get(ctx context.Context, key Key) (Object, error)
	// GetRange returns a result set of the specified key.
	GetRange(ctx context.Context, key Key, opts ...Option) (ResultSet, error)
	// Remove removes and returns the key-value object of the specified key.
	Remove(ctx context.Context, key Key) (Object, error)
	// RemoveRange removes the specified key-value objects.
	RemoveRange(ctx context.Context, key Key) error
}
