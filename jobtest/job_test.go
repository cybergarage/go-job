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
	"sync"
	"testing"

	"github.com/cybergarage/go-job/job"
)

func TestScheduleJobs(t *testing.T) {
	tests := []struct {
		kind string
		opts []any
		args []any
	}{
		{
			kind: "sum",
			opts: []any{
				job.WithExecutor(func(a, b int) int { return a + b }),
			},
			args: []any{1, 2},
		},
	}

	jobMgr := job.NewManager()

	if err := jobMgr.Start(); err != nil {
		t.Fatalf("Failed to start job manager: %v", err)
	}

	defer func() {
		if err := jobMgr.Stop(); err != nil {
			t.Errorf("Failed to stop job manager: %v", err)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			resHandler := func(job job.Instance, responses []any) {
				t.Logf("Job %s executed with responses: %v", job.Kind(), responses)
				wg.Done()
			}

			opts := append(
				tt.opts,
				job.WithKind(tt.kind),
				job.WithResponseHandler(resHandler),
			)
			j, err := job.NewJob(opts...)
			if err != nil {
				t.Fatalf("Failed to create job: %v", err)
			}
			err = jobMgr.RegisterJob(j)
			if err != nil {
				t.Fatalf("Failed to register job: %v", err)
			}
			err = jobMgr.ScheduleRegisteredJob(
				tt.kind,
				job.WithArguments(tt.args...))
			if err != nil {
				t.Fatalf("Failed to schedule job: %v", err)
			}

			wg.Wait()
		})
	}
}
