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

// Option represents a option.
type Option = any

// Store represents a store interface.
type Store interface {
	// Config returns the store configuration.
	Config
	// Name returns the name of the store.
	Name() string
	// Transact begin a new transaction.
	Transact(ctx context.Context, write bool) (Transaction, error)
	// Start starts the store.
	Start() error
	// Stop stops the store.
	Stop() error
}
