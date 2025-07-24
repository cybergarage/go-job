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
	"math/rand/v2"
	"time"
)

const (
	// NoRetry is a constant indicating that a job should not be retried.
	NoRetry = 0
	// RetryForever is a constant indicating that a job should retry indefinitely.
	RetryForever = -1

	// NoTimeout indicates no timeout limit.
	NoTimeout = 0
	// DefaultTimeout is the default timeout for jobs.
	DefaultTimeout = 0
)

// BackoffStrategy is a function type that defines how long to wait before retrying a job.
type BackoffStrategy func() time.Duration

// Policy defines the interface for job scheduling, supporting crontab expressions.
type Policy interface {
	// MaxRetries returns the maximum number of retries allowed for the job.
	MaxRetries() int
	// Priority returns the priority of the job.
	Priority() Priority
	// Timeout returns the timeout duration for the job.
	Timeout() time.Duration
	// Backoff returns the delay time before retrying a job.
	Backoff() time.Duration
	// Map returns a map representation of the job instance.
	Map() map[string]any
	// String returns a string representation of the job instance.
	String() string
}

// PolicyOption defines a function that configures a job policy.
type PolicyOption func(*policy)

// policy implements the JobPolicy interface using a crontab spec string.
type policy struct {
	maxRetries int
	priority   Priority
	timeout    time.Duration
	backoffFn  BackoffStrategy
}

// WithMaxRetries sets the maximum number of retries for the job policy.
func WithMaxRetries(count int) PolicyOption {
	return func(s *policy) {
		s.maxRetries = count
	}
}

// WithPriority sets the priority for the job policy.
func WithPriority(priority Priority) PolicyOption {
	return func(s *policy) {
		s.priority = priority
	}
}

// WithTimeout sets the timeout duration for the job policy.
func WithTimeout(duration time.Duration) PolicyOption {
	return func(s *policy) {
		s.timeout = duration
	}
}

// WithNoTimeout sets the job policy to have no timeout limit.
func WithNoTimeout() PolicyOption {
	return func(s *policy) {
		s.timeout = NoTimeout
	}
}

// WithHighPriority sets the job policy to high priority.
func WithHighPriority() PolicyOption {
	return func(s *policy) {
		s.priority = HighPriority
	}
}

// WithLowPriority sets the job policy to low priority.
func WithLowPriority() PolicyOption {
	return func(s *policy) {
		s.priority = LowPriority
	}
}

// WithInfiniteRetries sets the job policy to retry indefinitely.
func WithInfiniteRetries() PolicyOption {
	return func(s *policy) {
		s.maxRetries = RetryForever
	}
}

// WithBackoffStrategy sets the function to determine the delay before retrying a job.
func WithBackoffStrategy(fn BackoffStrategy) PolicyOption {
	return func(s *policy) {
		s.backoffFn = fn
	}
}

// WithBackoffDuration sets a fixed backoff duration with random jitter for the job policy.
func WithBackoffDuration(duration time.Duration) PolicyOption {
	return func(s *policy) {
		s.backoffFn = func() time.Duration {
			return time.Duration(float64(duration) * (0.8 + 0.4*(rand.Float64())))
		}
	}
}

func newPolicy() *policy {
	return &policy{
		maxRetries: NoRetry,         // Default to no retries
		priority:   DefaultPriority, // Default priority
		timeout:    DefaultTimeout,  // Default timeout
		backoffFn: func() time.Duration {
			return time.Duration(0)
		},
	}
}

// NewPolicy creates a new JobPolicy instance from a crontab spec string.
func NewPolicy(opts ...PolicyOption) (*policy, error) {
	js := newPolicy()
	for _, opt := range opts {
		opt(js)
	}
	return js, nil
}

// MaxRetries returns the maximum number of retries allowed for the job.
func (p *policy) MaxRetries() int {
	return p.maxRetries
}

// Priority returns the priority of the job.
func (p *policy) Priority() Priority {
	return p.priority
}

// Timeout returns the timeout duration for the job.
func (p *policy) Timeout() time.Duration {
	return p.timeout
}

// Backoff returns the delay time before retrying a job.
func (p *policy) Backoff() time.Duration {
	if p.backoffFn == nil {
		return 0
	}
	return p.backoffFn()
}

// Map returns a map representation of the job instance.
func (p *policy) Map() map[string]any {
	m := map[string]any{
		maxRetriesKey: p.MaxRetries(),
		priorityKey:   p.Priority(),
		timeoutKey:    p.Timeout().String(),
	}
	return m
}

// String returns a string representation of the job instance.
func (p *policy) String() string {
	return fmt.Sprintf("%v", p.Map())
}
