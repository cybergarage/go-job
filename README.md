# go-job

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/go-job)
[![test](https://github.com/cybergarage/go-job/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/go-job/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-job.svg)](https://pkg.go.dev/github.com/cybergarage/go-job)
 [![Go Report Card](https://img.shields.io/badge/go%20report-A%2B-brightgreen)](https://goreportcard.com/report/github.com/cybergarage/go-job) 
 [![codecov](https://codecov.io/gh/cybergarage/go-job/graph/badge.svg?token=OCU5V0H3OX)](https://codecov.io/gh/cybergarage/go-job)

`go-job` is a flexible and extensible job management library for Go.  
It provides job creation, registration, scheduling, execution, monitoring, and state/history management for distributed and local systems.

## Features

- Job creation with custom executors and arguments
- Job registration and deregistration
- FIFO job queue and scheduling
- Job instance management and state transitions
- Job execution with response and error handlers
- Persistent and pluggable job store interface
- Job state history logging and querying
- Fault tolerance and retry mechanisms

## Getting Started

### Installation

```sh
go get github.com/cybergarage/go-job
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/cybergarage/go-job/job"
)

func main() {
	// Create a job manager
	mgr := NewManager()

	// Register a job with a custom executor
	sumJob, _ := NewJob(
		WithKind("sum"),
		WithExecutor(func(a, b int) int { return a + b }),
		WithScheduleAt(time.Now()),
		WithResponseHandler(func(inst Instance, res []any) {
			fmt.Printf("Result: %v\n", res)
		}),
	)

	// Register the job with the manager
	mgr.RegisterJob(sumJob)

	// Schedule the registered job
	mgr.ScheduleRegisteredJob("sum", WithArguments(1, 2))

	// Start the job manager
	mgr.Start()

	// Wait for the job to complete
	mgr.StopWithWait()
}
```

## Features

### Job Execution with Arbitrary Functions

`go-job` allows you to register and execute jobs with any function signature using Go's `any` type for arguments and results. This flexibility enables you to define custom executors that accept and return any data types, making it easy to integrate with various workflows.

For example, you can register a job that adds two numbers, concatenates strings, or performs any custom logic:

```go
sumJob, _ := NewJob(
	WithKind("sum"),
	WithExecutor(func(a, b int) int { return a + b }),
)
type concatOpt struct {
		a string
		b string
	}
concatJob, _ := NewJob(
	WithKind("concat"),
	WithExecutor(func(opt concatOpt) string { return opt.a + " "+ opt.b }),
)
```

When scheduling jobs, you can pass arguments of any type, and the executor will receive them as parameters. The results are also handled as `any`, allowing flexible response handling.

```go
mgr.ScheduleJob(sumJob, WithArguments(sumOpt{1, 2}))
mgr.ScheduleJob(concatJob, WithArguments(concatOpt{"Hello", "world!"}))
```

This design makes `go-job` suitable for a wide range of use cases, from simple arithmetic to complex business logic.

### Flexible Scheduling

`go-job` provides flexible scheduling options to fit various job execution requirements. You can schedule jobs to run immediately, at a specific time, or after a certain delay. The library supports both one-time and recurring job executions.

#### Schedule at a Specific Time

You can schedule a job to run at a particular time using `WithScheduleAt`:

```go
mgr.ScheduleJob(job, WithScheduleAt(time.Now().Add(10 * time.Minute)))
```

#### Schedule with a Delay

To run a job after a delay, use `WithScheduleAfter`:

```go
mgr.ScheduleJob(job, WithScheduleAfter(5 * time.Second))
```

#### Cron-style Scheduling

You can schedule jobs using a cron expression with `WithCrontabSpec`. This allows you to define complex recurring schedules similar to Unix cron syntax:

```go
mgr.ScheduleJob(job, WithCrontabSpec("0 0 * * *")) // Runs every day at midnight
```

This feature is useful for jobs that need to run on specific days, times, or intervals defined by a crontab specification.

### Flexible Queue Priority and Worker Management

`go-job` supports prioritizing jobs in the queue and managing multiple workers for concurrent job execution. You can assign priorities to jobs, ensuring that higher-priority jobs are executed before lower-priority ones.

#### Job Priority

When creating or scheduling a job, use `WithPriority` to set its priority. Higher values indicate higher priority:

```go
mgr.ScheduleJob(job, WithPriority(10))
```

Jobs with higher priority are dequeued and executed before those with lower priority, allowing you to control the order of job processing.

#### Worker Pool Management

You can configure the number of worker goroutines that process jobs concurrently. This enables efficient resource utilization and parallel job execution:

```go
mgr, _ := NewManager(
	WithNumWorkers(5), // Set the number of workers to 5
)
mgr.Start()
mgr.ResizeWorkers(10) // Dynamically resize the worker pool to 10 workers
```

This flexibility allows you to scale job processing based on your application's needs, balancing throughput and resource usage.