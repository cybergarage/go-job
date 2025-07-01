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

// JobSchedule defines the interface for job scheduling, supporting crontab expressions.
type JobSchedule interface {
	// Spec returns the crontab spec string.
	Spec() string
	// Next returns the next scheduled time after the given time.
	Next(time.Time) time.Time
}

// JobScheduleOption defines a function that configures a job schedule.
type JobScheduleOption func(*jobSchedule) error

// jobSchedule implements the JobSchedule interface using a crontab spec string.
type jobSchedule struct {
	spec     string
	schedule cron.Schedule
}

// WithCronSpec sets the crontab spec string for the job schedule.
func WithCronSpec(spec string) JobScheduleOption {
	return func(js *jobSchedule) error {
		js.spec = spec
		var err error
		js.schedule, err = cron.ParseStandard(spec)
		if err != nil {
			return err
		}
		return nil
	}
}

// NewJobSchedule creates a new JobSchedule instance from a crontab spec string.
func NewJobSchedule(opts ...JobScheduleOption) (*jobSchedule, error) {
	js := &jobSchedule{
		spec:     "",
		schedule: nil, // Default to nil, must be set by options
	}
	for _, opt := range opts {
		if err := opt(js); err != nil {
			return nil, err
		}
	}
	return js, nil
}

// Spec returns the crontab spec string for the job schedule.
func (js *jobSchedule) Spec() string {
	return js.spec
}

// Next returns the next scheduled time after the given time.
func (js *jobSchedule) Next(t time.Time) time.Time {
	if js.schedule == nil {
		return t // If no schedule is set, return the given time
	}
	return js.schedule.Next(t)
}
