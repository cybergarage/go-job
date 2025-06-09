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

// CounterJob is a simple job that maintains a counter.
type CounterJob struct {
	Value int
}

// NewCounterJob creates a new instance of CounterJob with an initial value of 0.
func NewCounterJob() *CounterJob {
	return &CounterJob{Value: 0}
}

// Job that increments the counter
func (c *CounterJob) IncrementJob() func() error {
	return func() error {
		c.Value++
		return nil
	}
}
