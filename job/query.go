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

// Queue is an interface that defines queue operations for job management.
type Quey interface {
	JobState() JobState
}

type query struct {
	jobState JobState
}

// QueryOption is a function that configures a job query.
type QueryOption func(*query)

// WithQueryJobState sets the job state for the query.
func WithQueryJobState(state JobState) QueryOption {
	return func(q *query) {
		q.jobState = state
	}
}

// NewQuery creates a new instance of the job query.
func NewQuery(opts ...QueryOption) *query {
	q := &query{
		jobState: JobCreated, // Default state is Created
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// JobState returns the job state of the query.
func (q *query) JobState() JobState {
	return q.jobState
}
