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

package jobtest

import (
	"testing"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kv/memdb"
)

func StoreTest(t *testing.T, store kv.Store) {
	t.Helper()

	if err := store.Start(); err != nil {
		t.Fatalf("failed to start store: %v", err)
	}
	defer func() {
		if err := store.Stop(); err != nil {
			t.Fatalf("failed to stop store: %v", err)
		}
	}()
}

func TestStores(t *testing.T) {
	stores := []kv.Store{
		memdb.NewStore(),
	}

	for _, store := range stores {
		t.Run(store.Name(), func(t *testing.T) {
			StoreTest(t, store)
		})
	}
}
