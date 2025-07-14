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

package job

import "fmt"

// Priority represents the priority of a job. The lower the number, the higher the priority.
type Priority int

const (
	// HighPriority is the high priority for jobs.
	HighPriority = 0
	// MediumPriority is the medium priority for jobs.
	MediumPriority = 5
	// DefaultPriority is the default priority for jobs.
	DefaultPriority = MediumPriority
	// LowPriority is the low priority for jobs.
	LowPriority = 10
)

// Equal checks if the priority is equal to another priority.
func (p Priority) Equal(other Priority) bool {
	return p == other
}

// Lower checks if the priority is lower than another priority.
func (p Priority) Lower(other Priority) bool {
	return p > other
}

// Higher checks if the priority is higher than another priority.
func (p Priority) Higher(other Priority) bool {
	return p < other
}

// String returns the string representation of the priority.
func (p Priority) String() string {
	switch p {
	case HighPriority:
		return "High"
	case MediumPriority:
		return "Medium"
	case LowPriority:
		return "Low"
	default:
		return fmt.Sprintf("%d", p)
	}
}
