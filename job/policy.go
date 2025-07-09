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

// Policy defines the interface for job scheduling, supporting crontab expressions.
type Policy interface {
	// MaxRetries returns the maximum number of retries allowed for the job.
	MaxRetries() int
	// Priority returns the priority of the job.
	Priority() int
	// Timeout returns the timeout duration for the job.
	Timeout() time.Duration
}

// PolicyOption defines a function that configures a job policy.
type PolicyOption func(*policy) error

// policy implements the JobPolicy interface using a crontab spec string.
type policy struct {
	maxRetries int
	priority   int
	timeout    time.Duration
}

// WithMaxRetries sets the maximum number of retries for the job policy.
func WithMaxRetries(count int) PolicyOption {
	return func(s *policy) error {
		s.maxRetries = count
		return nil
	}
}

// WithPriority sets the priority for the job policy.
func WithPriority(priority int) PolicyOption {
	return func(s *policy) error {
		s.priority = priority
		return nil
	}
}

// WithTimeout sets the timeout duration for the job policy.
func WithTimeout(duration time.Duration) PolicyOption {
	return func(s *policy) error {
		s.timeout = duration
		return nil
	}
}

// WithNoTimeout sets the job policy to have no timeout limit.
func WithNoTimeout() PolicyOption {
	return func(s *policy) error {
		s.timeout = NoTimeout
		return nil
	}
}

// WithHighPriority sets the job policy to high priority.
func WithHighPriority() PolicyOption {
	return func(s *policy) error {
		s.priority = HighPriority
		return nil
	}
}

// WithLowPriority sets the job policy to low priority.
func WithLowPriority() PolicyOption {
	return func(s *policy) error {
		s.priority = LowPriority
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
		maxRetries: NoRetry,         // Default to no retries
		priority:   DefaultPriority, // Default priority
		timeout:    DefaultTimeout,  // Default timeout
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
