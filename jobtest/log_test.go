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

func TestLogLevelFilter(t *testing.T) {
	tests := []struct {
		name     string
		level    job.LogLevel
		contains job.LogLevel
		expected bool
	}{
		{"Contains Info", job.LogAll, job.LogInfo, true},
		{"Contains Error", job.LogAll, job.LogError, true},
		{"Contains Warning", job.LogAll, job.LogWarn, true},
		{"Does not contain Info", job.LogError, job.LogInfo, false},
		{"Does not contain Error", job.LogWarn, job.LogError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Contains(tt.contains); got != tt.expected {
				t.Errorf("Contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}
