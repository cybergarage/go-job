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

package job

import (
	"context"
)

// Store is the interface for job instance storage.
type Store interface {
	PendingStore
}

// PendingStore is an interface that defines methods for managing job instances.
type PendingStore interface {
	// EnqueueInstance stores a job instance in the store.
	EnqueueInstance(ctx context.Context, job Instance) error
	// RemoveInstance removes a job instance from the store by its unique identifier.
	RemoveInstance(ctx context.Context, job Instance) error
	// ListInstances lists all job instances in the store.
	ListInstances(ctx context.Context) ([]Instance, error)
}

// HistoryStore is an interface that defines methods for managing job instance state history.
type HistoryStore interface {
	// AddInstanceRecord adds a new state record for a job instance.
	AddInstanceRecord(ctx context.Context, job Instance, record InstanceRecord) error
	// ListInstanceRecords lists all state records for a job instance.
	ListInstanceRecords(ctx context.Context, job Instance) ([]InstanceRecord, error)
}
