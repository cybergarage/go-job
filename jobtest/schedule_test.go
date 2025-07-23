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

package jobtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/go-job/job"
)

func TestSchedules(t *testing.T) {
	tests := []struct {
		cronSpec   string
		scheduleAt time.Time
	}{
		{
			cronSpec:   "",
			scheduleAt: time.Now(),
		},
		{
			cronSpec:   "",
			scheduleAt: time.Now().Add(10 * time.Hour),
		},
		{
			cronSpec:   "0 0 * * *", // every day at midnight
			scheduleAt: time.Now(),
		},
		{
			cronSpec:   "*/5 * * * *", // every 5 minutes
			scheduleAt: time.Time{},
		},
		{
			cronSpec:   "0 9 * * 1-5", // 9am on weekdays
			scheduleAt: time.Time{},
		},
		{
			cronSpec:   "15 14 1 * *", // 2:15pm on the first of every month
			scheduleAt: time.Time{},
		},
		{
			cronSpec:   "0 0 29 2 *", // leap day (Feb 29)
			scheduleAt: time.Time{},
		},
		{
			cronSpec:   "0 0 1 1 *", // New Year's Day
			scheduleAt: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q/%v", tt.cronSpec, tt.scheduleAt), func(t *testing.T) {
			sched, err := job.NewSchedule(
				job.WithCrontabSpec(tt.cronSpec),
				job.WithScheduleAt(tt.scheduleAt),
			)
			if err != nil {
				t.Fatalf("Failed to create schedule: %v", err)
			}
			scheduleAt := sched.Next()
			if scheduleAt.IsZero() {
				t.Errorf("Expected a valid next schedule time, got zero value")
			}
		})
	}
}
