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

func TestJob(t *testing.T) {
	jobMgr := job.NewManager()

	if err := jobMgr.Start(); err != nil {
		t.Fatalf("Failed to start job manager: %v", err)
	}

	if err := jobMgr.Stop(); err != nil {
		t.Fatalf("Failed to stop job manager: %v", err)
	}
}

func TestNewJob(t *testing.T) {
	tests := []struct {
		opts []any
	}{
		{
			opts: []any{
				job.WithExecutor(func(a, b int) int { return a + b }),
				job.WithArguments(1, 2),
			},
		},
	}

	for _, tt := range tests {
		_, err := job.NewJob(tt.opts...)
		if err != nil {
			t.Fatalf("Failed to create job: %v", err)
		}
	}
}
