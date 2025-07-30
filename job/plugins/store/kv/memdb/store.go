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

package memdb

import (
	"github.com/cybergarage/go-job/job/plugins/store/kv"
)

// Store represents a Memdb store service instance.
type Store struct {
	kv.Config
	*Database
}

// NewStore returns a new memdb store instance.
func NewStore() kv.Store {
	return &Store{
		Config: kv.NewConfig(
			kv.WithUniqueKeys(true), // default to unique keys
		),
		Database: nil,
	}
}

// Name returns the name of this memdb store.
func (store *Store) Name() string {
	return "memdb"
}

// Start starts this memdb.
func (store *Store) Start() error {
	db, err := NewDatabase()
	if err != nil {
		return err
	}
	store.Database = db
	return nil
}

// Stop stops this memdb.
func (store *Store) Stop() error {
	return nil
}
