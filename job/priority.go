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

import (
	"fmt"
	"strconv"
)

// Priority represents the priority of a job. A lower value means a higher priority, similar to the Unix nice value.
type Priority int

const (
	// HighPriority is the high priority for jobs.
	HighPriority = Priority(0)
	// MediumPriority is the medium priority for jobs.
	MediumPriority = Priority(5)
	// DefaultPriority is the default priority for jobs.
	DefaultPriority = MediumPriority
	// LowPriority is the low priority for jobs.
	LowPriority = Priority(10)
)

// NewPriorityFrom creates a Priority from various input types.
func NewPriorityFrom(a any) (Priority, error) {
	switch v := a.(type) {
	case int:
		return Priority(v), nil
	case string:
		p, err := strconv.Atoi(v)
		if err != nil {
			return DefaultPriority, err
		}
		return Priority(p), nil
	default:
		return DefaultPriority, fmt.Errorf("invalid priority value: %v", a)
	}
}

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
		return strconv.Itoa(int(p))
	}
}
