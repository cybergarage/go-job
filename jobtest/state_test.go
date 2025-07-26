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

func TestJobState(t *testing.T) {

	tests := []struct {
		state job.JobState
	}{
		{state: job.JobCreated},
		{state: job.JobScheduled},
		{state: job.JobProcessing},
		{state: job.JobCancelled},
		{state: job.JobTimedOut},
		{state: job.JobCompleted},
		{state: job.JobTerminated},
	}

	for _, tt := range tests {
		t.Run("JobState "+tt.state.String(), func(t *testing.T) {
			if !tt.state.Is(tt.state) {
				t.Errorf("expected JobState %s to be equal to itself", tt.state.String())
			}
			if !tt.state.Is(job.JobStateAll) {
				t.Errorf("expected JobState %s to be part of JobStateAll", tt.state.String())
			}
			stateStr := tt.state.String()
			if stateStr == "" {
				t.Errorf("expected JobState %s to have a string representation", tt.state)
			}
			parsedState, err := job.NewStateFrom(stateStr)
			if err != nil {
				t.Errorf("failed to parse JobState from string: %v", err)
			}
			if parsedState != tt.state {
				t.Errorf("expected JobState %s to equal parsed state %s", tt.state.String(), parsedState.String())
			}
			parsedStateStr := parsedState.String()
			if parsedStateStr == "" {
				t.Errorf("expected parsed JobState %s to have a string representation", parsedState.String())
			}
			if parsedStateStr != stateStr {
				t.Errorf("expected parsed JobState string %s to equal original string %s", parsedStateStr, stateStr)
			}
		})
	}
}
