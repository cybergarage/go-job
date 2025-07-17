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

func TestPriorityCompares(t *testing.T) {
	tests := []struct {
		name           string
		p1             job.Priority
		p2             job.Priority
		eqExpected     bool
		lowerExpected  bool
		higherExpected bool
	}{
		{
			name:           "High vs High",
			p1:             job.HighPriority,
			p2:             job.HighPriority,
			eqExpected:     true,
			lowerExpected:  false,
			higherExpected: false,
		},
		{
			name:           "High vs Medium",
			p1:             job.HighPriority,
			p2:             job.MediumPriority,
			eqExpected:     false,
			lowerExpected:  false,
			higherExpected: true,
		},
		{
			name:           "High vs Low",
			p1:             job.HighPriority,
			p2:             job.LowPriority,
			eqExpected:     false,
			lowerExpected:  false,
			higherExpected: true,
		},
		{
			name:           "Medium vs High",
			p1:             job.MediumPriority,
			p2:             job.HighPriority,
			eqExpected:     false,
			lowerExpected:  true,
			higherExpected: false,
		},
		{
			name:           "Medium vs Medium",
			p1:             job.MediumPriority,
			p2:             job.MediumPriority,
			eqExpected:     true,
			lowerExpected:  false,
			higherExpected: false,
		},
		{
			name:           "Medium vs Low",
			p1:             job.MediumPriority,
			p2:             job.LowPriority,
			eqExpected:     false,
			lowerExpected:  false,
			higherExpected: true,
		},
		{
			name:           "Low vs High",
			p1:             job.LowPriority,
			p2:             job.HighPriority,
			eqExpected:     false,
			lowerExpected:  true,
			higherExpected: false,
		},
		{
			name:           "Low vs Medium",
			p1:             job.LowPriority,
			p2:             job.MediumPriority,
			eqExpected:     false,
			lowerExpected:  true,
			higherExpected: false,
		},
		{
			name:           "Low vs Low",
			p1:             job.LowPriority,
			p2:             job.LowPriority,
			eqExpected:     true,
			lowerExpected:  false,
			higherExpected: false,
		},
		{
			name:           "Custom vs High",
			p1:             job.Priority(3),
			p2:             job.HighPriority,
			eqExpected:     false,
			lowerExpected:  true,
			higherExpected: false,
		},
		{
			name:           "Custom vs Low",
			p1:             job.Priority(7),
			p2:             job.LowPriority,
			eqExpected:     false,
			lowerExpected:  false,
			higherExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p1.Equal(tt.p2); got != tt.eqExpected {
				t.Errorf("Priority.Equal() = %v, want %v", got, tt.eqExpected)
			}
			if got := tt.p1.Lower(tt.p2); got != tt.lowerExpected {
				t.Errorf("Priority.Lower() = %v, want %v", got, tt.lowerExpected)
			}
			if got := tt.p1.Higher(tt.p2); got != tt.higherExpected {
				t.Errorf("Priority.Higher() = %v, want %v", got, tt.higherExpected)
			}
			if tt.eqExpected && tt.p1.String() != tt.p2.String() {
				t.Errorf("Priority.String() mismatch: %v != %v", tt.p1.String(), tt.p2.String())
			}
		})
	}
}
