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

// Store is an interface that defines methods for storing and retrieving jobs.
// It provides a way to persist job data, allowing for job management
// across application restarts or failures.
// The Store interface is designed to be implemented by various storage
// backends, such as databases, in-memory stores, database systems.
// Implementations of this interface should handle the serialization
// and deserialization of job data, as well as any necessary locking
// or concurrency control to ensure data integrity.
type Store interface {
	// StoreJob stores a job in the store.
	StoreJob(ctx context.Context, job Job) error
	// UpdateJob updates an existing job in the store.
	UpdateJob(ctx context.Context, job Job) error
	// RemoveJob removes a job from the store by its unique identifier.
	RemoveJob(ctx context.Context, job Job) error
	// ListJobs lists all jobs in the store.
	ListJobs(ctx context.Context) ([]Job, error)
	// ClearJobs clears all jobs from the store.
	ClearJobs(ctx context.Context) error
}
