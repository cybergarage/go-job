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
	"time"

	"github.com/robfig/cron/v3"
)

const (
	// RetryForever is a constant indicating that a job should retry indefinitely.
	RetryForever = -1
)

// Schedule defines the interface for job scheduling, supporting crontab expressions.
type Schedule interface {
	// CrontabSpec returns the crontab spec string.
	CrontabSpec() string
	// Next returns the next scheduled time.
	Next() time.Time
}

// ScheduleOption defines a function that configures a job schedule.
type ScheduleOption func(*schedule) error

// schedule implements the JobSchedule interface using a crontab spec string.
type schedule struct {
	crontabSpec  string
	cronSchedule cron.Schedule
	scheduleAt   time.Time
	maxRetries   int
	retryCount   int // Number of retries attempted
}

// WithCrontabSpec sets the crontab spec string for the job schedule.
func WithCrontabSpec(spec string) ScheduleOption {
	return func(js *schedule) error {
		if len(spec) == 0 {
			js.cronSchedule = nil
			return nil
		}
		js.crontabSpec = spec
		var err error
		js.cronSchedule, err = cron.ParseStandard(spec)
		if err != nil {
			return err
		}
		return nil
	}
}

// WithSchedule sets the cron.Schedule for the job schedule.
func WithScheduleAt(t time.Time) ScheduleOption {
	return func(js *schedule) error {
		js.scheduleAt = t
		return nil
	}
}

// WithMaxRetries sets the maximum number of retries for the job schedule.
func WithMaxRetries(count int) ScheduleOption {
	return func(s *schedule) error {
		s.maxRetries = count
		return nil
	}
}

// WithInfiniteRetries sets the job schedule to retry indefinitely.
func WithInfiniteRetries() ScheduleOption {
	return func(s *schedule) error {
		s.maxRetries = RetryForever
		return nil
	}
}

func newSchedule() *schedule {
	return &schedule{
		crontabSpec:  "",
		cronSchedule: nil,
		scheduleAt:   time.Now(),
		maxRetries:   0, // Default to no retries
		retryCount:   0,
	}
}

// NewSchedule creates a new JobSchedule instance from a crontab spec string.
func NewSchedule(opts ...ScheduleOption) (*schedule, error) {
	js := newSchedule()
	for _, opt := range opts {
		if err := opt(js); err != nil {
			return nil, err
		}
	}
	return js, nil
}

// CrontabSpec returns the crontab spec string for the job schedule.
func (js *schedule) CrontabSpec() string {
	return js.crontabSpec
}

// Next returns the next scheduled time.
func (js *schedule) Next() time.Time {
	if js.cronSchedule != nil {
		return js.cronSchedule.Next(time.Now())
	}
	return js.scheduleAt
}
