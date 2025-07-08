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
	InstanceStore
}

// InstanceStore is an interface that defines methods for managing job instances.
type InstanceStore interface {
	// AddInstance stores a job instance in the store.
	AddInstance(ctx context.Context, job Instance) error
	// RemoveInstance removes a job instance from the store by its unique identifier.
	RemoveInstance(ctx context.Context, job Instance) error
	// ListInstances lists all job instances in the store.
	ListInstances(ctx context.Context) ([]Instance, error)
	// ListInstancesByState lists all job instances in the store by their state.
	ListInstancesByState(ctx context.Context, state JobState) ([]Instance, error)
}

// HistoryStore is an interface that defines methods for managing job instance state history.
type HistoryStore interface {
	AddInstanceHistory(ctx context.Context, job Instance, state JobState) error
	// ListInstanceHistory lists all state history records for a job instance.
	ListInstanceHistory(ctx context.Context, job Instance) ([]JobState, error)
}
