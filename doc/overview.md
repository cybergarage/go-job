# Feature Overview and Usage Guide

`go-job` is a flexible and extensible job scheduling and execution library for Go that supports arbitrary function execution, custom scheduling, job monitoring, priority queuing, and distributed operation.

<figure>
<img src="img/job-framework.png" alt="job framework" />
</figure>

This document provides a comprehensive overview of the features and usage of `go-job`.

## Features

`go-job` provides:

- Arbitrary function registration

- Rich scheduling options

- Prioritized and scalable execution

- Strong observability

- Pluggable, distributed storage

Use it to build robust, scalable job systems in Go.

### Arbitrary Function Execution

`go-job` allows registration and execution of **any** function using Go’s `any` type for arguments and results. The executor can be defined with any number of input and output parameters or with struct definitions.

#### Simple Function Example

A job with no input parameters and no return value can be defined as follows:

    job, err := NewJob(
        WithKind("hello (no input and no return)"),
        WithExecutor(func()  {
            fmt.Println("Hello, world!")
        }),
    )

Then schedule this job with no arguments simply by:

    mgr.ScheduleJob(job)

#### Function with Arguments Example

A job with two input parameters and no return value can be defined like this:

    job, err := NewJob(
        WithKind("sum (two input and no output)"),
        WithExecutor(func(x int, y int) {
            fmt.Println(x + y)
        }),
    )

Then schedule jobs with arguments:

    mgr.ScheduleJob(job, WithArguments(42, 58))

#### Function with Arguments and Result Example

A job with two input parameters and one output can be defined like this:

    job, err := NewJob(
        WithKind("concat (two input and one output)"),
        WithExecutor(func(a string, b string) string {
            return a + ", " + b
        }),
        WithCompleteProcessor(func(ji Instance, res []any) {
            // In this case, log the result to the go-job manager
            ji.Infof("%v", res[0])
        }),
    )

Use `WithCompleteProcessor()` to capture the result of a job execution. This is useful when the job has a return value.

Then schedule jobs with arguments:

    mgr.ScheduleJob(job, WithArguments("Hello", "world"))

#### Function with Struct Input and Output

A job with one struct input and one struct output can be defined like this:

    type concatString struct {
        a string
        b string
        s string
    }

    job, err := NewJob(
        WithKind("concat (one struct input and one struct output)"),
        WithExecutor(func(param *concatString) *concatString {
            // Store the concatenated string result in the input struct, and return it
            param.s = param.a + ", " + param.b
            return param
        }),
        WithCompleteProcessor(func(ji Instance, res []any) {
            // In this case, log the result to the go-job manager
            ji.Infof("%v", res[0])
        }),
    )

In this case, the result is also stored in the struct field `s`.

Then schedule the jobs with arguments by:

    mgr.ScheduleJob(job, WithArguments(&concatString{"Hello", "world!"}))

This approach supports diverse function signatures and is ideal for both simple and complex use cases. For additional examples, see the [Examples](https://pkg.go.dev/github.com/cybergarage/go-job/job#NewJob) section in the [![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-job.svg)](https://pkg.go.dev/github.com/cybergarage/go-job).

### Flexible Scheduling

Schedule jobs:

- **Immediately** (default)

- **At a specific time**

- **After a delay**

- **On a recurring cron schedule**

#### Schedule at a Specific Time

    mgr.ScheduleJob(job, WithScheduleAt(time.Now().Add(10 * time.Minute)))

#### Delay Execution

    mgr.ScheduleJob(job, WithScheduleAfter(5 * time.Second))

#### Cron Scheduling

    mgr.ScheduleJob(job, WithCrontabSpec("0 0 * * *")) // daily at midnight

Supports standard cron format: `min hour dom month dow`.

### Queue Priority & Worker Management

#### Job Priority

    mgr.ScheduleJob(job, WithPriority(0)) // high-priority

Higher-priority jobs are executed before lower-priority ones.

#### Dynamic Worker Pool

    mgr, _ := NewManager(WithNumWorkers(5))
    mgr.Start()
    mgr.ResizeWorkers(10)

Allows concurrent execution and real-time scalability.

### Job Observation

#### Handlers for Response and Error

    job, err := NewJob(
        WithKind("observe"),
        WithExecutor(func(x int) int { return x * 2 }),
        WithCompleteProcessor(func(inst Instance, res []any) {
            inst.Infof("Result: %v", res)
        }),
        WithTerminateProcessor(func(inst Instance, err error) {
            inst.Errorf("Error: %v", err)
        }),
    )
    mgr.ScheduleJob(job, WithArguments(42))

#### List All job Instances

##### List All Queued and Executed Job Instances

       query := job.NewQuery() // queries all job instances (any state)
        jis, err := mgr.LookupInstances(query)
        if err != nil {
            t.Errorf("Failed to lookup job instance: %v", err)
        }
        for _, ji := range jis {
            fmt.Printf("Job Instance: %s, UUID: %s, State: %s\n", ji.Kind(), ji.UUID(), ji.State())
        }

##### List Terminated Job Instances

        query := job.NewQuery(
            job.WithQueryKind("sum"), // filter by job kind
            job.WithQueryState(job.JobTerminated), // filter by terminated state
        )
        jis, err := mgr.LookupInstances(query)
        if err != nil {
            t.Errorf("Failed to lookup job instance: %v", err)
        }
        for _, ji := range jis {
            fmt.Printf("Job Instance: %s, State: %s\n", ji.Kind(), ji.State())
        }

#### Retrieve History and Logs for Job Instances

##### State History

    states := mgr.ProcessHistory(ji)
    for _, s := range states {
        fmt.Printf("State: %s at %v\n", s.State(), s.Timestamp())
    }

#### Log History

    logs := mgr.ProcessLogs(ji)
    for _, log := range logs {
        fmt.Printf("[%s] %v: %s\n", log.Level(), log.Timestamp(), log.Message())
    }

Provides auditability and debugging capability for each job instance.

#### Distributed Support via Store Interface

`go-job` supports pluggable storage via the `Store` interface.

    distStore := NewMyDistributedStore(...)
    mgr, _ := NewManager(WithStore(distStore))

By implementing a custom store (e.g., etcd, FoundationDB), job metadata and execution state can be shared across nodes.

This enables:

- Distributed scheduling

- Cross-node job coordination

- State persistence across restarts

- Fault-tolerant execution

To learn more about the `Store` interface, see [Design and Architecture](design.md) and [Extension Guide ](plugin-guide.md) documentation.
