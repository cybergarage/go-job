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
	"strconv"
	"testing"

	"github.com/cybergarage/go-job/job"
)

func TestResizeWorkers(t *testing.T) {
	mgr, err := job.NewManager()
	if err != nil {
		t.Fatalf("failed to create job manager: %v", err)
	}

	if n := mgr.NumWorkers(); n != job.DefaultWorkerNum {
		t.Fatalf("expected %d worker, got %d", job.DefaultWorkerNum, n)
	}

	tests := []struct {
		newSize  int
		expected bool
	}{
		{newSize: 0, expected: false},
		{newSize: -1, expected: false},
		{newSize: 3, expected: true},
		{newSize: 5, expected: true},
		{newSize: 2, expected: true},
	}

	for _, tt := range tests {
		t.Run("Resize to "+strconv.Itoa(tt.newSize), func(t *testing.T) {
			err := mgr.ResizeWorkers(tt.newSize)
			if !tt.expected {
				if err == nil {
					t.Errorf("expected error when resizing to %d workers, but got none", tt.newSize)
				}
				return
			}
			if err != nil {
				t.Errorf("failed to resize workers: %v", err)
				return
			}
			if n := mgr.NumWorkers(); n != tt.newSize {
				t.Errorf("expected %d workers after resizing, got %d", tt.newSize, n)
			}
		})
	}
}
