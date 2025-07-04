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
	// ScheduleJob schedules a job instance with the given job and options.
	// It creates a new job instance and enqueues it in the job queue.
	ScheduleJob(job Job, opts ...any) error
}

// SchedulerOption is a function that configures a job scheduler.
type SchedulerOption func(*scheduler)

// WithSchedulerQueue sets the job queue for the scheduler.
func WithSchedulerQueue(queue Queue) SchedulerOption {
	return func(s *scheduler) {
		s.Queue = queue
	}
}

type scheduler struct {
	Queue
}

// NewScheduler creates a new instance of Scheduler.
func NewScheduler(opts ...SchedulerOption) *scheduler {
	s := &scheduler{
		Queue: NewQueue(WithQueueStore(NewMemStore())),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ScheduleJob schedules a job instance with the given job and options.
// It creates a new job instance and enqueues it in the job queue.
func (s *scheduler) ScheduleJob(job Job, opts ...any) error {
	opts = append(opts,
		WithExecutor(job.Handler().Execute()),
		WithErrorHandler(job.Handler().ErrorHandler()),
		WithResponseHandler(job.Handler().ResponseHandler()),
	)
	ji, err := NewInstance(opts...)
	if err != nil {
		return err
	}
	if err := s.Queue.Enqueue(ji); err != nil {
		return err
	}
	return nil
}
