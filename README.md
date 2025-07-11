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
    mgr := job.NewManager()
    mgr.Start()
    defer mgr.Stop()

    // Register a job with a custom executor
    sumJob, _ := job.NewJob(
        job.WithKind("sum"),
        job.WithExecutor(func(a, b int) int { return a + b }),
        job.WithResponseHandler(func(inst job.Instance, res []any) {
            fmt.Printf("Result: %v\n", res)
        }),
    )
    mgr.RegisterJob(sumJob)

    // Schedule the registered job
    mgr.ScheduleRegisteredJob("sum", job.WithArguments(1, 2))
}
```

### Job Lifecycle

- **Job**: Defines the task logic and metadata.
- **Instance**: Represents a scheduled or running job with its own state and arguments.
- **State**: Each instance tracks its state (e.g., Scheduled, Running, Completed, Failed, Canceled).
- **Store**: Pluggable interface for FIFO job queue and state history persistence.

### Interfaces

#### Job

- `WithKind`, `WithExecutor`, `WithArguments`, `WithResponseHandler` for flexible job definition.

#### Manager

- `RegisterJob`, `UnregisterJob`, `ScheduleRegisteredJob`, `Start`, `Stop`.

#### Store

- `EnqueueInstance`, `DequeueInstance`, `RemoveInstance`, `ListInstances`
- `LogInstanceRecord`, `ListInstanceRecords`

#### Instance

- `Job()`, `Kind()`, `UUID()`, `ScheduledAt()`, `State()`, `UpdateState()`, `Process()`

### Advanced

- Custom job stores can be implemented by satisfying the `Store` interface.
- Job state transitions and history are tracked for monitoring and debugging.
- Supports distributed and local execution models.

## License

Apache License 2.0