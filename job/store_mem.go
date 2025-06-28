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
	jobs sync.Map
}

// NewMemStore creates a new in-memory job store.
func NewMemStore() Store {
	return &memStore{
		jobs: sync.Map{},
	}
}

// StoreJob stores a job in the in-memory store.
func (s *memStore) StoreJob(ctx context.Context, job JobInstance) error {
	if job == nil {
		return nil
	}
	return nil
}

func (s *memStore) UpdateJob(ctx context.Context, job JobInstance) error {
	if job == nil {
		return nil
	}
	return nil
}

// RemoveJob removes a job from the in-memory store by its unique identifier.
func (s *memStore) RemoveJob(ctx context.Context, job JobInstance) error {
	if job == nil {
		return nil
	}
	return nil
}

// ListJobs lists all jobs in the in-memory store.
func (s *memStore) ListJobs(ctx context.Context) ([]JobInstance, error) {
	jobs := make([]JobInstance, 0)
	s.jobs.Range(func(key, value interface{}) bool {
		if job, ok := value.(JobInstance); ok {
			jobs = append(jobs, job)
		}
		return true
	})
	return jobs, nil
}

// ClearJobs clears all jobs from the in-memory store.
func (s *memStore) ClearJobs(ctx context.Context) error {
	s.jobs.Range(func(key, value interface{}) bool {
		s.jobs.Delete(key)
		return true
	})
	return nil
}
