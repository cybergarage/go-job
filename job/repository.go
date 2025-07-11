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

// Repository is an interface that defines methods for managing job registrations and scheduling.
type Repository interface {
	// Registry defines the methods for managing job registrations.
	Registry
	// Scheduler defines the methods for scheduling jobs.
	Scheduler
}

// RepositoryOption is a function that configures a job repository.
type RepositoryOption func(*repository)

// WithRepositoryStore sets the store for the job repository.
func WithRepositoryStore(store Store) RepositoryOption {
	return func(r *repository) {
		r.store = store
	}
}

type repository struct {
	store Store
	Scheduler
	Registry
}

func NewRepository(opts ...RepositoryOption) *repository {
	repo := &repository{
		store:     NewLocalStore(),
		Scheduler: nil,
		Registry:  nil,
	}
	for _, opt := range opts {
		opt(repo)
	}
	queue := NewQueue(WithQueueStore(repo.store))
	repo.Scheduler = NewScheduler(WithSchedulerQueue(queue))

	repo.Registry = NewRegistry()

	return repo
}
