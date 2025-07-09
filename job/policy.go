package job

import "time"

const (
	// NoRetry is a constant indicating that a job should not be retried.
	NoRetry = 0
	// RetryForever is a constant indicating that a job should retry indefinitely.
	RetryForever = -1

	// DefaultPriority is the default priority for jobs.
	DefaultPriority = 0
	// HighPriority represents high priority jobs.
	HighPriority = 10
	// LowPriority represents low priority jobs.
	LowPriority = -10

	// NoTimeout indicates no timeout limit.
	NoTimeout = 0
	// DefaultTimeout is the default timeout for jobs.
	DefaultTimeout = 30 * time.Minute
)

// RetryDelay is a function type that defines how long to wait before retrying a job.
type RetryDelay func() time.Duration

// Policy defines the interface for job scheduling, supporting crontab expressions.
type Policy interface {
	// MaxRetries returns the maximum number of retries allowed for the job.
	MaxRetries() int
	// Priority returns the priority of the job.
	Priority() int
	// Timeout returns the timeout duration for the job.
	Timeout() time.Duration
	// RetryDelay returns the delay time before retrying a job.
	RetrayDelay() time.Duration
}

// PolicyOption defines a function that configures a job policy.
type PolicyOption func(*policy)

// policy implements the JobPolicy interface using a crontab spec string.
type policy struct {
	maxRetries   int
	priority     int
	timeout      time.Duration
	retryDelayFn RetryDelay
}

// WithMaxRetries sets the maximum number of retries for the job policy.
func WithMaxRetries(count int) PolicyOption {
	return func(s *policy) {
		s.maxRetries = count
	}
}

// WithPriority sets the priority for the job policy.
func WithPriority(priority int) PolicyOption {
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

// WithRetryDelay sets the function to determine the delay before retrying a job.
func WithRetryDelay(fn RetryDelay) PolicyOption {
	return func(s *policy) {
		s.retryDelayFn = fn
	}
}

// WithRetryDelayDuration sets a fixed delay duration before retrying a job.
func WithRetryDelayDuration(duration time.Duration) PolicyOption {
	return func(s *policy) {
		s.retryDelayFn = func() time.Duration {
			return duration
		}
	}
}

func newPolicy() *policy {
	return &policy{
		maxRetries: NoRetry,         // Default to no retries
		priority:   DefaultPriority, // Default priority
		timeout:    DefaultTimeout,  // Default timeout
		retryDelayFn: func() time.Duration {
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
func (p *policy) Priority() int {
	return p.priority
}

// Timeout returns the timeout duration for the job.
func (p *policy) Timeout() time.Duration {
	return p.timeout
}

// RetryDelay returns the delay time before retrying a job.
func (p *policy) RetryDelay() time.Duration {
	if p.retryDelayFn == nil {
		return 0
	}
	return p.retryDelayFn()
}
