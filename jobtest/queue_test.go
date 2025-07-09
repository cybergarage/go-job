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

	"github.com/cybergarage/go-job/job"
)

func QueueStoreTest(t *testing.T, store job.Store) {
	t.Helper()

	jobs := []job.Instance{}

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
			t.Errorf("Failed to dequeue job: %v", err)
			return
		}
		if lastJob == nil {
			lastJob = job
			continue
		}
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
