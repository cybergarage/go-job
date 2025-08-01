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
	"context"
	"fmt"
	"testing"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kv/memdb"
	"github.com/cybergarage/go-job/job/plugins/store/kvutil"
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

	// Set / Get tests

	setterTests := []struct {
		key kv.Key
		val []byte
	}{
		{
			key: kv.Key("test1"),
			val: []byte("value1"),
		},
		{
			key: kv.Key("test2"),
			val: []byte("value2"),
		},
		{
			key: kv.Key("test3"),
			val: []byte("value3"),
		},
	}

	for _, test := range setterTests {
		t.Run("Set "+test.key.String(), func(t *testing.T) {
			tx, err := store.Transact(context.Background(), true)
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Cancel(context.Background())

			obj := kv.NewObject(test.key, test.val)
			if err := tx.Set(context.Background(), obj); err != nil {
				t.Fatalf("failed to set object: %v", err)
			}

			if err := tx.Commit(context.Background()); err != nil {
				t.Fatalf("failed to commit transaction: %v", err)
			}

			tx, err = store.Transact(context.Background(), false)
			if err != nil {
				t.Fatalf("failed to begin read transaction: %v", err)
			}
			defer tx.Cancel(context.Background())

			retrievedObj, err := tx.Get(context.Background(), test.key)
			if err != nil {
				t.Fatalf("failed to get object: %v", err)
			}

			if !retrievedObj.Equal(obj) {
				t.Errorf("expected %v, got %v", obj, retrievedObj)
			}
		})
		t.Run("Remove "+test.key.String(), func(t *testing.T) {
			tx, err := store.Transact(context.Background(), true)
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Cancel(context.Background())

			if err := tx.Remove(context.Background(), test.key); err != nil {
				t.Fatalf("failed to remove object: %v", err)
			}

			if err := tx.Commit(context.Background()); err != nil {
				t.Fatalf("failed to commit transaction: %v", err)
			}

			tx, err = store.Transact(context.Background(), false)
			if err != nil {
				t.Fatalf("failed to begin read transaction: %v", err)
			}
			defer tx.Cancel(context.Background())

			_, err = tx.Get(context.Background(), test.key)
			if err == nil {
				t.Errorf("expected error when getting removed object, got nil")
			}
		})
	}

	// Set / GetRange tests

	rangeKey := kv.Key("range1")
	rangeTests := []struct {
		key kv.Key
		val []byte
	}{
		{
			key: rangeKey,
			val: []byte("value1"),
		},
		{
			key: rangeKey,
			val: []byte("value2"),
		},
		{
			key: rangeKey,
			val: []byte("value3"),
		},
	}

	keys := make([]kv.Key, len(rangeTests))
	objs := make([]kv.Object, len(rangeTests))
	for i, test := range rangeTests {
		keys[i] = kv.Key(fmt.Sprintf("%s%d", test.key, i))
		objs[i] = kv.NewObject(keys[i], test.val)
		t.Run(fmt.Sprintf("SetRange %s %s", test.key.String(), string(test.val)), func(t *testing.T) {
			tx, err := store.Transact(context.Background(), true)
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Cancel(context.Background())

			if err := tx.Set(context.Background(), objs[i]); err != nil {
				t.Fatalf("failed to set object: %v", err)
			}

			if err := tx.Commit(context.Background()); err != nil {
				t.Fatalf("failed to commit transaction: %v", err)
			}
		})
		t.Run(fmt.Sprintf("GetRange %s %s", test.key.String(), string(test.val)), func(t *testing.T) {
			tx, err := store.Transact(context.Background(), false)
			if err != nil {
				t.Fatalf("failed to begin read transaction: %v", err)
			}
			defer tx.Cancel(context.Background())

			rs, err := tx.GetRange(context.Background(), test.key)
			if err != nil {
				t.Fatalf("failed to get range: %v", err)
			}

			retrievedObjs, err := kvutil.ReadAll(rs)
			if err != nil {
				t.Fatalf("failed to read all objects from range: %v", err)
			}
			if len(retrievedObjs) != i+1 {
				t.Fatalf("expected %d objects, got %d", i+1, len(retrievedObjs))
			}

			for _, obj := range retrievedObjs {
				found := false
				for _, expectedObj := range objs[:i+1] {
					if obj.Equal(expectedObj) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("unexpected object: %v", obj)
				}
			}
		})
	}

	t.Run("RemoveRange "+rangeKey.String(), func(t *testing.T) {
		tx, err := store.Transact(context.Background(), true)
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}
		defer tx.Cancel(context.Background())

		if err := tx.RemoveRange(context.Background(), rangeKey); err != nil {
			t.Fatalf("failed to remove object: %v", err)
		}

		if err := tx.Commit(context.Background()); err != nil {
			t.Fatalf("failed to commit transaction: %v", err)
		}

		tx, err = store.Transact(context.Background(), false)
		if err != nil {
			t.Fatalf("failed to begin read transaction: %v", err)
		}
		defer tx.Cancel(context.Background())

		rs, err := tx.GetRange(context.Background(), rangeKey)
		if err != nil {
			t.Fatalf("expected error when getting removed object, got nil")
		}

		retrievedObjs, err := kvutil.ReadAll(rs)
		if err != nil {
			t.Fatalf("failed to read all objects from range: %v", err)
		}
		if 0 < len(retrievedObjs) {
			t.Errorf("expected no objects after removal, got %d", len(retrievedObjs))
		}
	})
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
