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
	"github.com/cybergarage/go-job/job/plugins/store"
	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kvutil"
	// "github.com/cybergarage/go-job/job/plugins/store/kv/valkey"
)

func InstanceQueueStoreTest(t *testing.T, store job.Store) {
	t.Helper()

	if err := store.Start(); err != nil {
		t.Errorf("Failed to start store: %v", err)
		return
	}

	defer func() {
		if err := store.Stop(); err != nil {
			t.Errorf("Failed to stop store: %v", err)
			return
		}
	}()

	if err := store.Clear(); err != nil {
		t.Errorf("Failed to clear store: %v", err)
		return
	}

	now := time.Now()

	tests := []struct {
		opts []any
	}{
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(-10 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.HighPriority),
				job.WithScheduleAt(now.Add(-10 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(-5 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(-1 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.DefaultPriority),
				job.WithScheduleAt(now.Add(-1 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(-500 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.HighPriority),
				job.WithScheduleAt(now.Add(-500 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(-2 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.DefaultPriority),
				job.WithScheduleAt(now.Add(-2 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(-100 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.HighPriority),
				job.WithScheduleAt(now.Add(-3 * time.Hour)),
			},
		},
		{
			opts: []any{
				job.WithPriority(job.LowPriority),
				job.WithScheduleAt(now.Add(-700 * time.Hour)),
			},
		},
	}

	ctx := t.Context()

	// Enqueue jobs

	jobs := []job.Instance{}
	for _, tt := range tests {
		job, err := job.NewInstance(tt.opts...)
		if err != nil {
			t.Errorf("Failed to create job: %v", err)
			return
		}
		jobs = append(jobs, job)
	}

	queue := job.NewInstanceQueue(job.WithInstanceQueueStore(store))

	for _, job := range jobs {
		if err := queue.Enqueue(ctx, job); err != nil {
			t.Errorf("Failed to enqueue job: %v", err)
			return
		}
	}

	switch store := store.(type) {
	case kv.Store:
		objs, err := store.Dump(t.Context())
		if err != nil {
			t.Errorf("Failed to dump store: %v", err)
		}
		kvutil.LogObjects(t, objs)
	}

	// Dequeue jobs

	var lastJob job.Instance
	for range jobs {
		job, err := queue.Dequeue(ctx)
		if err != nil {
			t.Fatalf("Failed to dequeue job: %v", err)
			return
		}
		t.Logf("Dequeued job: %d, scheduled at: %s", job.Policy().Priority(), job.ScheduledAt().Format(time.RFC3339))
		if lastJob == nil {
			lastJob = job
			continue
		}

		lp := lastJob.Policy().Priority()
		cp := job.Policy().Priority()
		switch {
		case cp.Higher(lp):
			t.Fatalf("Job dequeued out of order: %d after %d", cp, lp)
		case cp.Equal(lp):
			if job.ScheduledAt().Before(lastJob.ScheduledAt()) {
				t.Fatalf("Job dequeued out of order: %s before %s", job.ScheduledAt().String(), lastJob.ScheduledAt().String())
			}
		}

		if job.Before(lastJob) {
			t.Fatalf("Job dequeued out of order: %v before %v", job.ScheduledAt(), lastJob.ScheduledAt())
		}

		lastJob = job
	}
}

func TestInstanceQueue(t *testing.T) {
	stores := []job.Store{
		job.NewLocalStore(),
		store.NewMemdbStore(),
		// store.NewValkeyStore(valkey.NewStoreOption()),
	}

	for _, store := range stores {
		t.Run(store.Name(), func(t *testing.T) {
			InstanceQueueStoreTest(t, store)
		})
	}
}
