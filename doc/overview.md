# Feature Overview and Usage Guide
---
generator: Asciidoctor 2.0.23
lang: en
title: Feature Overview and Usage Guide
viewport: width=device-width, initial-scale=1.0
---

<div id="header">

# Feature Overview and Usage Guide

</div>

<div id="content">

<div id="preamble">

<div class="sectionbody">

<div class="paragraph">

`go-job` is a flexible and extensible job scheduling and execution library for Go that supports arbitrary function execution, custom scheduling, job monitoring, priority queuing, and distributed operation.

</div>

<div class="imageblock">

<div class="content">

![job framework](img/job-framework.png)

</div>

</div>

<div class="paragraph">

This document provides a comprehensive overview of the features and usage of `go-job`.

</div>

</div>

</div>

<div class="sect1">

## Features

<div class="sectionbody">

<div class="paragraph">

`go-job` provides:

</div>

<div class="ulist">

- Arbitrary function registration

- Rich scheduling options

- Strong observability

- Prioritized and scalable execution

- Pluggable, distributed storage

</div>

<div class="paragraph">

Use it to build robust, scalable job systems in Go.

</div>

<div class="sect2">

### Arbitrary Function Execution

<div class="paragraph">

`go-job` allows you to register and execute **any Go function** as a job. You can use functions with different signatures - from simple functions with no parameters to complex functions with multiple inputs and outputs.

</div>

<div class="listingblock">

<div class="content">

``` highlight
// Executor can be any function type
type Executor any

// Examples of valid executors:
// func()                           // no input, no output
// func(int, string) bool           // multiple inputs, one output
// func(*MyStruct) error            // struct input, error output
// func(a, b int) (int, error)      // multiple inputs and outputs
```

</div>

</div>

<div class="paragraph">

This flexibility means you can:

</div>

<div class="ulist">

- Use functions with any number of parameters

- Return single values or multiple values

- Work with primitive types (int, string, bool)

- Pass complex structs as arguments

- Handle errors in your job functions

</div>

<div class="paragraph">

The `any` type allows `go-job` to work with your existing functions without requiring special interfaces or wrapper code.

</div>

<div class="sect3">

#### Simple Function Example

<div class="paragraph">

A job with no input parameters and no return value can be defined as follows:

</div>

<div class="listingblock">

<div class="content">

``` highlight
job, err := NewJob(
    WithKind("hello (no input and no return)"),
    WithExecutor(func()  {
        fmt.Println("Hello, world!")
    }),
)
```

</div>

</div>

<div class="paragraph">

Then schedule this job with no arguments simply by:

</div>

<div class="listingblock">

<div class="content">

``` highlight
mgr.ScheduleJob(job)
```

</div>

</div>

</div>

<div class="sect3">

#### Function with Arguments Example

<div class="paragraph">

A job with two input parameters and no return value can be defined like this:

</div>

<div class="listingblock">

<div class="content">

``` highlight
job, err := NewJob(
    WithKind("sum (two input and no output)"),
    WithExecutor(func(x int, y int) {
        fmt.Println(x + y)
    }),
)
```

</div>

</div>

<div class="paragraph">

Then schedule jobs with arguments:

</div>

<div class="listingblock">

<div class="content">

``` highlight
mgr.ScheduleJob(job, WithArguments(42, 58))
```

</div>

</div>

</div>

<div class="sect3">

#### Function with Arguments and Result Example

<div class="paragraph">

A job with two input parameters and one output can be defined like this:

</div>

<div class="listingblock">

<div class="content">

``` highlight
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
```

</div>

</div>

<div class="paragraph">

Use `WithCompleteProcessor()` to capture the result of a job execution. This is useful when the job has a return value.

</div>

<div class="paragraph">

Then schedule jobs with arguments:

</div>

<div class="listingblock">

<div class="content">

``` highlight
mgr.ScheduleJob(job, WithArguments("Hello", "world"))
```

</div>

</div>

</div>

<div class="sect3">

#### Function with Struct Input and Output

<div class="paragraph">

A job with one struct input and one struct output can be defined like this:

</div>

<div class="listingblock">

<div class="content">

``` highlight
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
```

</div>

</div>

<div class="paragraph">

In this case, the result is also stored in the struct field `s`.

</div>

<div class="paragraph">

Then schedule the jobs with arguments by:

</div>

<div class="listingblock">

<div class="content">

``` highlight
arg := &concatString{
    a: "Hello",
    b: "world!",
    s: "",
}
mgr.ScheduleJob(job, WithArguments(arg))
```

</div>

</div>

<div class="paragraph">

This approach supports diverse function signatures and is ideal for both simple and complex use cases. For additional examples, see the [Examples](https://pkg.go.dev/github.com/cybergarage/go-job/job#NewJob) section in the [<span class="image">![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-job.svg)</span>](https://pkg.go.dev/github.com/cybergarage/go-job).

</div>

</div>

</div>

<div class="sect2">

### Job Scheduling

<div class="paragraph">

`go-job` provides flexible scheduling options to run jobs when you need them:

</div>

<div class="ulist">

- **Immediately** - Jobs start executing right away (default behavior)

- **At a specific time** - Schedule jobs to run at an exact date and time

- **After a delay** - Wait a specified duration before starting execution

- **On a recurring schedule** - Use cron expressions for repeated execution

</div>

<div class="sect3">

#### Execute Jobs Immediately

<div class="paragraph">

By default, jobs are scheduled for immediate execution:

</div>

<div class="listingblock">

<div class="content">

``` highlight
// Runs immediately
mgr.ScheduleJob(job)
```

</div>

</div>

</div>

<div class="sect3">

#### Schedule at a Specific Time

<div class="paragraph">

Set an exact time for job execution:

</div>

<div class="listingblock">

<div class="content">

``` highlight
// Run 10 minutes from now
futureTime := time.Now().Add(10 * time.Minute)
mgr.ScheduleJob(job, WithScheduleAt(futureTime))

// Run at a specific date and time
specificTime := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
mgr.ScheduleJob(job, WithScheduleAt(specificTime))
```

</div>

</div>

</div>

<div class="sect3">

#### Delay Execution

<div class="paragraph">

Add a delay before the job starts:

</div>

<div class="listingblock">

<div class="content">

``` highlight
// Wait 5 seconds before execution
mgr.ScheduleJob(job, WithScheduleAfter(5 * time.Second))

// Wait 2 hours before execution
mgr.ScheduleJob(job, WithScheduleAfter(2 * time.Hour))
```

</div>

</div>

</div>

<div class="sect3">

#### Recurring Cron Scheduling

<div class="paragraph">

Use cron expressions for repeated job execution:

</div>

<div class="listingblock">

<div class="content">

``` highlight
// Run daily at midnight
mgr.ScheduleJob(job, WithCrontabSpec("0 0 * * *"))

// Run every weekday at 9 AM
mgr.ScheduleJob(job, WithCrontabSpec("0 9 * * 1-5"))

// Run every 30 minutes
mgr.ScheduleJob(job, WithCrontabSpec("*/30 * * * *"))
```

</div>

</div>

<div class="paragraph">

Cron format: `minute hour day-of-month month day-of-week`

</div>

</div>

</div>

<div class="sect2">

### Job Monitoring and Observability

<div class="paragraph">

`go-job` provides comprehensive monitoring capabilities to track job execution and understand system behavior. You can monitor jobs in real-time using event handlers, or query historical data using manager methods.

</div>

<div class="sect3">

#### Real-time Monitoring with Event Handlers

<div class="paragraph">

Monitor job execution as it happens by registering event handlers that respond to completion, termination, and state changes.

</div>

<div class="sect4">

##### Completion and Termination Handlers

<div class="paragraph">

Use `WithCompleteProcessor()` and `WithTerminateProcessor()` to handle successful completion and error termination:

</div>

<div class="listingblock">

<div class="content">

``` highlight
job, err := NewJob(
    ....,
    WithCompleteProcessor(func(inst Instance, res []any) {
        inst.Infof("Result: %v", res)
    }),
    WithTerminateProcessor(func(inst Instance, err error) {
        inst.Errorf("Error: %v", err)
    }),
)
```

</div>

</div>

</div>

<div class="sect4">

##### State Change Monitoring

<div class="paragraph">

Use `WithStateChangeProcessor()` to track every state transition throughout a job’s lifecycle:

</div>

<div class="listingblock">

<div class="content">

``` highlight
job, err := NewJob(
    ....,
    WithStateChangeProcessor(func(inst Instance, state JobState) error {
        inst.Infof("State changed to: %v", state)
        return nil
    }),
)
```

</div>

</div>

<div class="paragraph">

For details on job state transitions, refer to [Design and Architecture](design.md).

</div>

</div>

</div>

<div class="sect3">

#### Historical Data Queries

<div class="paragraph">

Query job instances and their execution history using manager methods.

</div>

<div class="sect4">

##### List All job Instances

<div class="paragraph">

With `Manager::LookupInstances()`, you can retrieve any job instance—whether it is scheduled, in progress, or already executed.

</div>

<div class="sect5">

###### List All Queued and Executed Job Instances

<div class="listingblock">

<div class="content">

``` highlight
    query := job.NewQuery() // queries all job instances (any state)
    jis, err := mgr.LookupInstances(query)
    if err != nil {
        t.Errorf("Failed to lookup job instance: %v", err)
    }
    for _, ji := range jis {
        fmt.Printf("Job Instance: %s, UUID: %s, State: %s\n", ji.Kind(), ji.UUID(), ji.State())
    }
```

</div>

</div>

</div>

<div class="sect5">

###### List Terminated Job Instances

<div class="listingblock">

<div class="content">

``` highlight
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
```

</div>

</div>

</div>

</div>

<div class="sect4">

##### Retrieve History and Logs for Job Instances

<div class="paragraph">

You can use manager methods to access the processing history and logs of any specified job instance.

</div>

<div class="sect5">

###### State History

<div class="paragraph">

With `Manager::LookupInstanceHistory`, you can retrieve the state history for the specified job instance.

</div>

<div class="listingblock">

<div class="content">

``` highlight
states := mgr.LookupInstanceHistory(ji)
for _, s := range states {
    fmt.Printf("State: %s at %v\n", s.State(), s.Timestamp())
}
```

</div>

</div>

<div class="paragraph">

For details on job state transitions, refer to [Design and Architecture](design.md).

</div>

</div>

<div class="sect5">

###### Log History

<div class="paragraph">

With `Manager::LookupInstanceLogs`, you can retrieve the log history for the specified job instance.

</div>

<div class="listingblock">

<div class="content">

``` highlight
logs := mgr.LookupInstanceLogs(ji)
for _, log := range logs {
    fmt.Printf("[%s] %v: %s\n", log.Level(), log.Timestamp(), log.Message())
}
```

</div>

</div>

<div class="paragraph">

Provides auditability and debugging capability for each job instance.

</div>

</div>

</div>

</div>

</div>

<div class="sect2">

### Priority Management & Worker Scaling

<div class="paragraph">

`go-job` allows you to control job execution order through priorities and dynamically scale workers to handle varying workloads.

</div>

<div class="sect3">

#### Job Priority Control

<div class="paragraph">

Assign priorities to jobs to control their execution order. Higher priority jobs are executed before lower priority ones. The priority value is an integer where lower values indicate higher priority (similar to Unix nice values).

</div>

<div class="sect4">

##### Set Priority During Job Creation

<div class="listingblock">

<div class="content">

``` highlight
// High priority job (executed first)
highPriorityJob, err := NewJob(
    WithKind("urgent-task"),
    WithPriority(0), // lower number = higher priority like Unix nice values
    WithExecutor(func() { fmt.Println("Urgent task executing") }),
)

// Low priority job (executed later)
lowPriorityJob, err := NewJob(
    WithKind("background-task"),
    WithPriority(200), // higher number = lower priority like Unix nice values
    WithExecutor(func() { fmt.Println("Background task executing") }),
)
```

</div>

</div>

</div>

<div class="sect4">

##### Override Priority at Schedule Time

<div class="paragraph">

You can override a job’s default priority when scheduling:

</div>

<div class="listingblock">

<div class="content">

``` highlight
// Schedule with default priority
mgr.ScheduleJob(normalJob) // uses job's configured priority

// Schedule with custom priority (overrides job's default priority)
mgr.ScheduleJob(normalJob, WithPriority(200)) // make this instance low priority
```

</div>

</div>

</div>

</div>

<div class="sect3">

#### Dynamic Worker Pool Management

<div class="paragraph">

Scale your worker pool up or down based on workload demands without stopping the manager.

</div>

<div class="sect4">

##### Set Initial Worker Count

<div class="listingblock">

<div class="content">

``` highlight
// Start with 5 workers
mgr, err := NewManager(WithNumWorkers(5))
mgr.Start()
```

</div>

</div>

</div>

<div class="sect4">

##### Scale Workers Dynamically

<div class="listingblock">

<div class="content">

``` highlight
// Scale up during high load
mgr.ResizeWorkers(10) // increase to 10 workers

// Scale down during low load
mgr.ResizeWorkers(3)  // reduce to 3 workers

// Get current worker count
count := mgr.NumWorkers()
fmt.Printf("Current workers: %d\n", count)
```

</div>

</div>

</div>

<div class="sect4">

##### Real-world Scaling Example

<div class="listingblock">

<div class="content">

``` highlight
// Monitor queue size and scale accordingly
query := job.NewQuery(
    job.WithQueryState(job.JobScheduled), // filter by scheduled state
)
jobs, _ := mgr.LookupInstances(query)
queueSize := len(jobs)
currentWorkers := mgr.NumWorkers()
if queueSize > currentWorkers*2 {
    // Scale up if queue is getting too long
    mgr.ResizeWorkers(currentWorkers + 2)
} else if queueSize == 0 && currentWorkers > 2 {
    // Scale down if no jobs queued
    mgr.ResizeWorkers(currentWorkers - 1)
}
```

</div>

</div>

<div class="paragraph">

This enables efficient resource utilization and responsive performance under varying workloads.

</div>

</div>

</div>

</div>

<div class="sect2">

### Remote Management with gRPC API

<div class="paragraph">

`go-job` provides a comprehensive gRPC API for remote job management, enabling you to schedule, monitor, and control jobs from external systems or distributed environments. This allows seamless integration with microservices, orchestration platforms, and remote applications.

</div>

</div>

<div class="sect2">

### Remote Operation with gRPC API

<div class="paragraph">

`go-job` provides a gRPC API for remote job management, scheduling, and monitoring. This enables integration with external systems and remote orchestration. The gRPC API offers full programmatic access to all core `go-job` functionality:

</div>

<div class="ulist">

- Remote job scheduling with arguments and timing options

- Real-time job monitoring and status queries

- Dynamic worker pool management

- Cross-platform compatibility through protocol buffers

- Secure communication with authentication support

</div>

<div class="paragraph">

The gRPC API uses protobuf messages for job definitions, arguments, and results. For more details, see the [grpc.proto](grpc-api.md) definition.

</div>

<div class="sect3">

#### Command-Line Interface (jobctl)

<div class="paragraph">

`go-job` provides a command-line interface called [jobctl](./cmd/cli/jobctl.md) to interact with the gRPC API. The following methods are available:

</div>

<div class="ulist">

- `ScheduleJob` - Schedule a new job remotely with arguments and scheduling options

- `ListJobs` - List all registered jobs and their metadata

- `ListInstances` - Query job instances by kind, state, or time range

</div>

<div class="paragraph">

For more details, see the [Command-Line Interface (jobctl)](./cmd/cli/jobctl.md) documentation.

</div>

</div>

</div>

<div class="sect2">

### Distributed Support via Store Interface

<div class="paragraph">

`go-job` supports pluggable storage through the `Store` interface. The following component diagram shows how multiple `go-job` instances can share a single store.

</div>

<div class="imageblock">

<div class="content">

![job store](img/job-store.png)

</div>

</div>

<div class="paragraph">

By implementing a custom store (e.g., etcd, FoundationDB), job metadata and execution state can be shared across nodes.

</div>

<div class="paragraph">

This enables:

</div>

<div class="ulist">

- Distributed scheduling

- Cross-node job coordination

- State persistence across restarts

- Fault-tolerant execution

</div>

<div class="paragraph">

To learn more about the `Store` interface, see [Design and Architecture](design.md) and [Extension Guide](plugin-guide.md) documentation.

</div>

</div>

</div>

</div>

</div>

<div id="footer">

<div id="footer-text">

Last updated 2025-08-04 22:15:12 +0900

</div>

</div>
