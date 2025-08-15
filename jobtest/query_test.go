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
	"github.com/google/uuid"
)

func TestQuery(t *testing.T) {
	t.Helper()

	tests := []struct {
		opts           []job.QueryOption
		expectedUUID   bool
		expectedKind   bool
		expectedState  bool
		expectedLevel  bool
		expectedBefore bool
		expectedAfter  bool
	}{
		{
			opts:           []job.QueryOption{},
			expectedUUID:   false,
			expectedKind:   false,
			expectedState:  false,
			expectedLevel:  false,
			expectedBefore: false,
			expectedAfter:  false,
		},
		{
			opts: []job.QueryOption{
				job.WithQueryUUID(uuid.Nil),
				job.WithQueryKind(""),
				job.WithQueryState(job.JobStateUnset),
				job.WithQueryLogLevel(job.LogNone),
				job.WithQueryBefore(time.Time{}),
				job.WithQueryAfter(time.Time{}),
			},
			expectedUUID:   false,
			expectedKind:   false,
			expectedState:  false,
			expectedLevel:  false,
			expectedBefore: false,
			expectedAfter:  false,
		},
		{
			opts: []job.QueryOption{
				job.WithQueryUUID(uuid.New()),
				job.WithQueryKind("test"),
				job.WithQueryState(job.JobCreated),
				job.WithQueryLogLevel(job.LogInfo),
				job.WithQueryBefore(time.Now()),
				job.WithQueryAfter(time.Now()),
			},
			expectedUUID:   true,
			expectedKind:   true,
			expectedState:  true,
			expectedLevel:  true,
			expectedBefore: true,
			expectedAfter:  true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("uuid(%t),kind(%t),state(%t)", tt.expectedUUID, tt.expectedKind, tt.expectedState), func(t *testing.T) {
			query := job.NewQuery(tt.opts...)
			uuid, ok := query.UUID()
			if ok != tt.expectedUUID {
				t.Errorf("expected UUID presence: %v, got: %v", tt.expectedUUID, ok)
			}
			if ok && uuid.String() == "" {
				t.Error("expected non-nil UUID")
			}

			kind, ok := query.Kind()
			if ok != tt.expectedKind {
				t.Errorf("expected Kind presence: %v, got: %v", tt.expectedKind, ok)
			}
			if ok && kind == "" {
				t.Error("expected non-empty Kind")
			}

			state, ok := query.State()
			if ok != tt.expectedState {
				t.Errorf("expected State presence: %v, got: %v", tt.expectedState, ok)
			}
			if ok && state == job.JobStateUnset {
				t.Error("expected non-unset State")
			}

			level, ok := query.LogLevel()
			if ok != tt.expectedLevel {
				t.Errorf("expected LogLevel presence: %v, got: %v", tt.expectedLevel, ok)
			}
			if ok && level == job.LogNone {
				t.Error("expected non-none LogLevel")
			}

			before, ok := query.Before()
			if ok != tt.expectedBefore {
				t.Errorf("expected Before presence: %v, got: %v", tt.expectedBefore, ok)
			}
			if ok && before.IsZero() {
				t.Error("expected non-zero Before")
			}

			after, ok := query.After()
			if ok != tt.expectedAfter {
				t.Errorf("expected After presence: %v, got: %v", tt.expectedAfter, ok)
			}
			if ok && after.IsZero() {
				t.Error("expected non-zero After")
			}
		})
	}
}
