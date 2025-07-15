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
	WithExecutor(func(opt concatOpt) string { return opt.a + opt.b }),
)
```

When scheduling jobs, you can pass arguments of any type, and the executor will receive them as parameters. The results are also handled as `any`, allowing flexible response handling.

```go	
mgr.ScheduleJob(sumJob, WithArguments(1, 2))
mgr.ScheduleJob(concatJob, WithArguments("Hello, ", "world!"))
```

This design makes `go-job` suitable for a wide range of use cases, from simple arithmetic to complex business logic.
