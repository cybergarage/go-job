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

func TestFilter(t *testing.T) {
	t.Helper()

	tests := []struct {
		opts           []job.FilterOption
		expectedBefore bool
		expectedAfter  bool
	}{
		{
			opts:           []job.FilterOption{},
			expectedBefore: false,
			expectedAfter:  false,
		},
		{
			opts: []job.FilterOption{
				job.WithFilterBefore(time.Time{}),
				job.WithFilterAfter(time.Time{}),
			},
			expectedBefore: false,
			expectedAfter:  false,
		},
		{
			opts: []job.FilterOption{
				job.WithFilterBefore(time.Now()),
				job.WithFilterAfter(time.Now()),
			},
			expectedBefore: true,
			expectedAfter:  true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("before(%t)", tt.expectedBefore), func(t *testing.T) {
			query := job.NewFilter(tt.opts...)
			before, ok := query.Before()
			if ok != tt.expectedBefore {
				t.Errorf("expected Before() to return %t, got %t", tt.expectedBefore, ok)
			}
			if tt.expectedBefore && before.IsZero() {
				t.Errorf("expected Before() to return a non-zero time, got zero time")
			}
			after, ok := query.After()
			if ok != tt.expectedAfter {
				t.Errorf("expected After() to return %t, got %t", tt.expectedAfter, ok)
			}
			if tt.expectedAfter && after.IsZero() {
				t.Errorf("expected After() to return a non-zero time, got zero time")
			}
		})
	}
}
