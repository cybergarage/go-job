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
	"fmt"
)

// repository is an interface that defines methods for managing job registrations and scheduling.
type repository interface {
	// Registry defines the methods for managing job registrations.
	registry
	// Scheduler defines the methods for scheduling jobs.
	scheduler
	// History defines the methods for managing job history.
	History
	// Clear clears all job queues, history, and logs.
	Clear() error
}

// repositoryOption is a function that configures a job repository.
type repositoryOption func(*repositoryImpl)

// withRepositoryStore sets the store for the job repository.
func withRepositoryStore(store Store) repositoryOption {
	return func(r *repositoryImpl) {
		r.store = store
	}
}

type repositoryImpl struct {
	registry
	scheduler
	History

	store Store
}

// newRepository creates a new instance of Repository with the given options.
func newRepository(opts ...repositoryOption) *repositoryImpl {
	repo := &repositoryImpl{
		store:     NewLocalStore(),
		scheduler: nil,
		registry:  nil,
		History:   nil,
	}

	for _, opt := range opts {
		opt(repo)
	}

	repo.registry = newRegistry()
	repo.scheduler = newScheduler(withSchedulerStore(repo.store))
	repo.History = newHistory(withHistoryStore(repo.store))

	return repo
}

// Clear clears all job queues, history, and logs.
func (repo *repositoryImpl) Clear() error {
	registryCleaner := func(ctx context.Context) error {
		return repo.registry.Clear()
	}
	historyCleaner := func(ctx context.Context) error {
		return repo.ClearHistory(NewFilter())
	}
	logCleaner := func(ctx context.Context) error {
		return repo.ClearLogs(NewFilter())
	}
	clearners := []func(context.Context) error{
		registryCleaner,
		historyCleaner,
		logCleaner,
		repo.store.ClearInstances,
	}
	for _, clear := range clearners {
		if err := clear(context.Background()); err != nil {
			return fmt.Errorf("failed to clear job repository: %w", err)
		}
	}
	return nil
}
