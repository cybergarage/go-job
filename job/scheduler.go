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
	ScheduleJob(job Job) error
}

// SchedulerOption is a function that configures a job scheduler.
type SchedulerOption func(*scheduler)

func WithSchedulerQueue(queue Queue) SchedulerOption {
	return func(s *scheduler) {
		s.Queue = queue
	}
}

type scheduler struct {
	Queue
}

// NewScheduler creates a new instance of Scheduler.
func NewScheduler() *scheduler {
	return &scheduler{
		Queue: nil, // Initialize with a default queue or nil
	}
}

// ScheduleJob schedules a job for execution.
func (s *scheduler) ScheduleJob(job Job) error {
	// Implementation for scheduling the job
	return nil
}
