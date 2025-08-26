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
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// ScheduleOption defines a function that configures a job schedule.
type ScheduleOption func(*schedule) error

// JitterGenerator defines a function that returns a random jitter duration for job scheduling.
type JitterGenerator func() time.Duration

// Schedule defines the interface for job scheduling, supporting crontab expressions.
type Schedule interface {
	// CrontabSpec returns the crontab spec string.
	CrontabSpec() string
	// IsScheduled returns true if the schedule has timing configuration.
	IsScheduled() bool
	// IsRecurring checks if the job is recurring.
	IsRecurring() bool
	// Next returns the next scheduled time.
	Next() time.Time
	// Jitter returns the jitter duration for the job.
	Jitter() JitterGenerator
	// Map returns a map representation of the job.
	Map() map[string]any
	// String returns a string representation of the job.
	String() string
}

// schedule implements the JobSchedule interface using a crontab spec string.
type schedule struct {
	crontabSpec  string
	cronSchedule cron.Schedule
	scheduleAt   time.Time
	jitterFunc   JitterGenerator
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

// WithScheduleNow sets the job schedule to the current time.
func WithScheduleAfter(d time.Duration) ScheduleOption {
	return func(js *schedule) error {
		return WithScheduleAt(time.Now().Add(d))(js)
	}
}

// WithJitter sets the job schedule jitter function.
func WithJitter(jitterFunc JitterGenerator) ScheduleOption {
	return func(js *schedule) error {
		js.jitterFunc = jitterFunc
		return nil
	}
}

// newCrontabSpecFrom creates a crontab spec string from various input types.
func newCrontabSpecFrom(a any) (string, error) {
	switch v := a.(type) {
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid crontab spec value: %v", a)
	}
}

func newSchedule(opts ...ScheduleOption) (*schedule, error) {
	js := &schedule{
		crontabSpec:  "",
		cronSchedule: nil,
		scheduleAt:   time.Time{},
		jitterFunc:   func() time.Duration { return 0 },
	}
	for _, opt := range opts {
		if err := opt(js); err != nil {
			return nil, err
		}
	}
	return js, nil
}

// NewSchedule creates a new schedule instance with the provided options.
// Available options include WithCrontabSpec() for cron-based scheduling and WithScheduleAt() for one-time scheduling.
// If no options are provided, the current time is used as the default schedule.
func NewSchedule(opts ...ScheduleOption) (Schedule, error) {
	return newSchedule(opts...)
}

// CrontabSpec returns the crontab spec string for the job schedule.
func (js *schedule) CrontabSpec() string {
	return js.crontabSpec
}

// IsRecurring checks if the job is recurring based on the crontab spec.
func (js *schedule) IsRecurring() bool {
	return js.cronSchedule != nil
}

// IsScheduled returns true if the schedule has timing configuration.
func (js *schedule) IsScheduled() bool {
	return js.cronSchedule != nil || !js.scheduleAt.IsZero()
}

// Jitter returns the jitter duration for the job.
func (js *schedule) Jitter() JitterGenerator {
	return js.jitterFunc
}

// Next returns the next scheduled time.
func (js *schedule) Next() time.Time {
	jitter := time.Duration(0)
	if js.jitterFunc != nil {
		jitter = js.jitterFunc()
	}
	if js.cronSchedule != nil {
		return js.cronSchedule.Next(time.Now()).Add(jitter)
	}
	return js.scheduleAt.Add(jitter)
}

// Map returns a map representation of the job schedule.
func (js *schedule) Map() map[string]any {
	m := map[string]any{}
	if 0 < len(js.crontabSpec) {
		m[crontabKey] = js.crontabSpec
	}
	if !js.scheduleAt.IsZero() {
		m[scheduleAtKey] = NewTimestampFromTime(js.scheduleAt).String()
	}
	return m
}

// String returns a string representation of the job schedule.
func (js *schedule) String() string {
	return fmt.Sprintf("%v", js.Map())
}
