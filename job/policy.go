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

const (
	// RetryForever is a constant indicating that a job should retry indefinitely.
	RetryForever = -1
)

// Policy defines the interface for job scheduling, supporting crontab expressions.
type Policy interface {
}

// PolicyOption defines a function that configures a job policy.
type PolicyOption func(*policy) error

// policy implements the JobPolicy interface using a crontab spec string.
type policy struct {
	maxRetries int
}

// WithMaxRetries sets the maximum number of retries for the job policy.
func WithMaxRetries(count int) PolicyOption {
	return func(s *policy) error {
		s.maxRetries = count
		return nil
	}
}

// WithInfiniteRetries sets the job policy to retry indefinitely.
func WithInfiniteRetries() PolicyOption {
	return func(s *policy) error {
		s.maxRetries = RetryForever
		return nil
	}
}

func newPolicy() *policy {
	return &policy{
		maxRetries: 0, // Default to no retries
	}
}

// NewPolicy creates a new JobPolicy instance from a crontab spec string.
func NewPolicy(opts ...PolicyOption) (*policy, error) {
	js := newPolicy()
	for _, opt := range opts {
		if err := opt(js); err != nil {
			return nil, err
		}
	}
	return js, nil
}
