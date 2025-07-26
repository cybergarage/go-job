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

func TestTimestamps(t *testing.T) {
	// Test NewTimestamp
	t1 := job.NewTimestamp()
	if t1.String() == "" {
		t.Error("NewTimestamp returned an empty string")
	}

	// Test NewTimestampFromTime
	t2 := job.NewTimestampFromTime(t1.Time())
	if !t1.Equal(t2) {
		t.Errorf("NewTimestampFromTime returned a different time: got %s, want %s", t2.String(), t1.String())
	}

	// Test NewTimestampFromString
	t3, err := job.NewTimestampFromString(t1.String())
	if err != nil {
		t.Errorf("NewTimestampFromString failed: %v", err)
	}
	if !t1.Equal(t3) {
		t.Errorf("NewTimestampFromString returned a different time: got %s, want %s", t3.String(), t1.String())
	}
}
