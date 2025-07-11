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

// Scheduler is responsible for scheduling jobs.
type Scheduler interface {
	// Queue returns the job queue used by the scheduler.
	Queue() Queue
	// ScheduleJobInstance schedules a job instance with the given options.
	ScheduleJobInstance(job Instance, opts ...any) error
}

// SchedulerOption is a function that configures a job scheduler.
type SchedulerOption func(*scheduler)

// WithSchedulerQueue sets the job queue for the scheduler.
func WithSchedulerQueue(queue Queue) SchedulerOption {
	return func(s *scheduler) {
		s.queue = queue
	}
}

// WithSchedulerStore sets the store for the scheduler.
func WithSchedulerStore(store Store) SchedulerOption {
	return func(s *scheduler) {
		s.store = store
	}
}

type scheduler struct {
	store Store
	queue Queue
}

// NewScheduler creates a new instance of Scheduler.
func NewScheduler(opts ...SchedulerOption) *scheduler {
	s := &scheduler{
		store: NewLocalStore(),
		queue: nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	if s.queue == nil {
		s.queue = NewQueue(WithQueueStore(s.store))
	}
	return s
}

// Queue returns the job queue used by the scheduler.
func (s *scheduler) Queue() Queue {
	return s.queue
}

// ScheduleJobInstance schedules a job instance by adding it to the queue.
func (s *scheduler) ScheduleJobInstance(job Instance, opts ...any) error {
	if err := s.queue.Enqueue(job); err != nil {
		return err
	}
	return nil
}
