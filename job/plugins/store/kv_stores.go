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

package store

import (
	"github.com/cybergarage/go-job/job/plugins"
	"github.com/cybergarage/go-job/job/plugins/store/kv/etcd"
	"github.com/cybergarage/go-job/job/plugins/store/kv/memdb"
	"github.com/cybergarage/go-job/job/plugins/store/kv/valkey"
)

// NewMemdbStore creates a new in-memory key-value store instance.
func NewMemdbStore() plugins.Store {
	return NewKvStoreWith(memdb.NewStore())
}

// NewValkeyStore creates a new Valkey key-value store instance.
func NewValkeyStore(option valkey.StoreOption) plugins.Store {
	return NewKvStoreWith(valkey.NewStore(option))
}

// NewEtcdStore creates a new Etcd key-value store instance.
func NewEtcdStore(option etcd.StoreOption) plugins.Store {
	return NewKvStoreWith(etcd.NewStore(option))
}
