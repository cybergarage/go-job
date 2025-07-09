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
	"sync"
)

type memStore struct {
	jobs    sync.Map
	history []InstanceRecord
}

// NewMemStore creates a new in-memory job store.
func NewMemStore() Store {
	return &memStore{
		jobs:    sync.Map{},
		history: []InstanceRecord{},
	}
}

// EnqueueInstance stores a job instance in the store.
func (s *memStore) EnqueueInstance(ctx context.Context, job Instance) error {
	if job == nil {
		return nil
	}
	s.jobs.Store(job.UUID(), job)
	return nil
}

// RemoveInstance removes a job instance from the store by its unique identifier.
func (s *memStore) RemoveInstance(ctx context.Context, job Instance) error {
	if job == nil {
		return nil
	}
	s.jobs.Delete(job.UUID())
	return nil
}

// ListInstances lists all job instances in the store.
func (s *memStore) ListInstances(ctx context.Context) ([]Instance, error) {
	jobs := make([]Instance, 0)
	s.jobs.Range(func(key, value interface{}) bool {
		if job, ok := value.(Instance); ok {
			jobs = append(jobs, job)
		}
		return true
	})
	return jobs, nil
}

// LogInstanceRecord adds a new state record for a job instance.
func (s *memStore) LogInstanceRecord(ctx context.Context, job Instance, record InstanceRecord) error {
	if job == nil {
		return nil
	}
	s.history = append(s.history, record)
	return nil
}

// ListInstanceRecords lists all state records for a job instance.
func (s *memStore) ListInstanceRecords(ctx context.Context, job Instance) ([]InstanceRecord, error) {
	if job == nil {
		return nil, nil
	}
	var records []InstanceRecord
	for _, record := range s.history {
		if record.UUID() == job.UUID() {
			records = append(records, record)
		}
	}
	return records, nil
}
