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
	"testing"
	"time"

	"github.com/cybergarage/go-job/job"
)

func QueueStoreTest(t *testing.T, store job.Store) {
	t.Helper()

	now := time.Now()

	tests := []struct {
		opts []any
	}{
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(10 * time.Microsecond)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.HighPriority),
				job.WithScheduleAt(now.Add(10 * time.Microsecond)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(5 * time.Millisecond)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(1 * time.Second)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.DefaultPriority),
				job.WithScheduleAt(now.Add(1 * time.Second)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(500 * time.Microsecond)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.HighPriority),
				job.WithScheduleAt(now.Add(500 * time.Microsecond)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(2 * time.Second)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.DefaultPriority),
				job.WithScheduleAt(now.Add(2 * time.Second)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(100 * time.Millisecond)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.HighPriority),
				job.WithScheduleAt(now.Add(3 * time.Millisecond)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(700 * time.Microsecond)),
			},
		},
	}

	jobs := []job.Instance{}
	for _, tt := range tests {
		job, err := job.NewInstance(tt.opts...)
		if err != nil {
			t.Errorf("Failed to create job: %v", err)
			return
		}
		jobs = append(jobs, job)
	}

	q := job.NewQueue(job.WithQueueStore(store))

	for _, job := range jobs {
		if err := q.Enqueue(job); err != nil {
			t.Errorf("Failed to enqueue job: %v", err)
			return
		}
	}

	var lastJob job.Instance
	for i := 0; i < len(jobs); i++ {
		job, err := q.Dequeue()
		if err != nil {
			t.Fatalf("Failed to dequeue job: %v", err)
			return
		}
		t.Logf("Dequeued job: %d, scheduled at: %s", job.Policy().Priority(), job.ScheduledAt().String())
		if lastJob == nil {
			lastJob = job
			continue
		}
		if job.Policy().Priority() > lastJob.Policy().Priority() {
			t.Fatalf("Job dequeued out of order: %v before %v", job.ScheduledAt(), lastJob.ScheduledAt())
		}
		if job.ScheduledAt().Before(lastJob.ScheduledAt()) {
			t.Fatalf("Job dequeued out of order: %v before %v", job.ScheduledAt(), lastJob.ScheduledAt())
		}
		if job.ScheduledAt().Before(lastJob.ScheduledAt()) {
			t.Fatalf("Job dequeued out of order: %v before %v", job.ScheduledAt(), lastJob.ScheduledAt())
		}
		lastJob = job
	}
}

func TestQueue(t *testing.T) {
	stores := []job.Store{
		job.NewLocalStore(),
	}

	for _, store := range stores {
		t.Run(store.Name(), func(t *testing.T) {
			QueueStoreTest(t, store)
		})
	}
}
