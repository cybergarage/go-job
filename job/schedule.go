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

// Schedule defines the interface for job scheduling, supporting crontab expressions.
type Schedule interface {
	// Spec returns the crontab spec string.
	Spec() string
	// Next returns the next scheduled time.
	Next() time.Time
}

// ScheduleOption defines a function that configures a job schedule.
type ScheduleOption func(*schedule) error

// schedule implements the JobSchedule interface using a crontab spec string.
type schedule struct {
	crontabSpec string
	schedule    cron.Schedule
	scheduleAt  time.Time
}

// WithCrontabSpec sets the crontab spec string for the job schedule.
func WithCrontabSpec(spec string) ScheduleOption {
	return func(js *schedule) error {
		js.crontabSpec = spec
		var err error
		js.schedule, err = cron.ParseStandard(spec)
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

// NewJobSchedule creates a new JobSchedule instance from a crontab spec string.
func NewJobSchedule(opts ...ScheduleOption) (*schedule, error) {
	js := &schedule{
		crontabSpec: "",
		schedule:    nil, // Default to nil, must be set by options
		scheduleAt:  time.Now(),
	}
	for _, opt := range opts {
		if err := opt(js); err != nil {
			return nil, err
		}
	}
	return js, nil
}

// Spec returns the crontab spec string for the job schedule.
func (js *schedule) Spec() string {
	return js.crontabSpec
}

// Next returns the next scheduled time.
func (js *schedule) Next() time.Time {
	if js.schedule != nil {
		return js.schedule.Next(time.Now())
	}
	return js.scheduleAt
}
