<div id="header">

# Feature Overview and Usage Guide

</div>

<div id="content">

<div id="preamble">

<div class="sectionbody">

<div class="paragraph">

This document provides a comprehensive overview of the features and usage of `go-job`.

</div>

<div id="toc" class="toc">

<div id="toctitle" class="title">

Table of Contents:

</div>

- [Features](#_features)
  - [Arbitrary function registration](#_arbitrary_function_registration)
  - [Rich scheduling options](#_rich_scheduling_options)
  - [Strong observability](#_strong_observability)
  - [Prioritized and scalable execution](#_prioritized_and_scalable_execution)
  - [Pluggable, distributed storage](#_pluggable_distributed_storage)
- [Usage Guide](#_usage_guide)
  - [Arbitrary Function Execution](#_arbitrary_function_execution)
  - [Job Scheduling](#_job_scheduling)
  - [Job Monitoring and Observability](#_job_monitoring_and_observability)
  - [Priority Management & Worker Scaling](#_priority_management_worker_scaling)
  - [System Jobs](#_system_jobs)
  - [Remote Management with gRPC API](#_remote_management_with_grpc_api)
  - [Distributed Support via Store Interface](#_distributed_support_via_store_interface)

</div>

</div>

</div>

<div class="sect1">

## Features

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

`go-job` offers the following key features:

</div>

<div class="sect2">

### Arbitrary function registration

<div class="paragraph">

Register any Go function as a job, regardless of its signature. This allows you to schedule and execute a wide variety of tasks, from simple functions to complex business logic, without needing to conform to a specific interface.

</div>

</div>

<div class="sect2">

### Rich scheduling options

<div class="paragraph">

Schedule jobs to run immediately, at a specific time, after a delay, or on a recurring schedule using cron expressions. This flexibility enables you to automate tasks according to your application’s needs.

</div>

</div>

<div class="sect2">

### Strong observability

<div class="paragraph">

Monitor job execution in real time, track state transitions, and access detailed logs and history for each job instance. This provides transparency and makes it easy to debug and audit job processing.

</div>

</div>

<div class="sect2">

### Prioritized and scalable execution

<div class="paragraph">

Assign priorities to jobs to control execution order, and dynamically scale the number of worker processes to handle varying workloads efficiently. This ensures that urgent tasks are handled promptly and resources are used optimally.

</div>

</div>

<div class="sect2">

### Pluggable, distributed storage

<div class="paragraph">

Use a variety of storage backends (such as Valkey, etcd, or in-memory) to persist job metadata and state. Distributed storage support enables robust, fault-tolerant job processing across multiple nodes.

</div>

</div>

</div>

</div>

<div class="sect1">

## Usage Guide

<div class="sectionbody">

<div class="paragraph">

This section provides practical guidance on how to use `go-job` in your applications. Follow these examples and explanations to quickly integrate `go-job` into your workflow and take full advantage of its features.

</div>

<div class="sect2">

### Arbitrary Function Execution

<div class="paragraph">

With `go-job`, you can register and execute **any Go function** as a job—no matter its signature. This means you aren’t limited to a specific interface or function type: you can use anything from simple, no-argument functions to complex functions with multiple parameters and return values.

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
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

Thanks to this flexibility, you can:

</div>

<div class="ulist">

- Use functions with any number and type of parameters

- Return single or multiple values

- Work with basic types (int, string, bool) or complex structs

- Handle errors directly in your job functions

</div>

<div class="paragraph">

Because `go-job` uses the Go `any` type for executors, you don’t need to write adapters or wrappers—just pass your existing functions as they are.

</div>

<div class="paragraph">

This approach makes it easy to use `go-job` for both simple tasks and advanced workflows. For more examples, see the [Examples](https://pkg.go.dev/github.com/cybergarage/go-job/job#NewJob) section in the [<span class="image">![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/go-job.svg)</span>](https://pkg.go.dev/github.com/cybergarage/go-job).

</div>

<div class="sect3">

#### Simple Function Example

<div class="paragraph">

A job with no input parameters and no return value can be defined as follows:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
job, err := job.NewJob(
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

``` CodeRay
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

``` CodeRay
job, err := job.NewJob(
    WithKind("sum (two input and no output)"),
    WithExecutor(func(x int, y int) {
        fmt.Println(x + y)
    }),
)
```

</div>

</div>

<div class="paragraph">

You can schedule jobs by passing arguments of any type. For example, to schedule a job with integer arguments:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithArguments(42, 58))
```

</div>

</div>

<div class="paragraph">

Or, to schedule a job with string arguments:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithArguments("42", "58"))
```

</div>

</div>

<div class="paragraph">

`go-job` will automatically convert the provided arguments to match the types expected by the job function, so you can use the most convenient format for your use case.

</div>

</div>

<div class="sect3">

#### Function with Arguments and Result Example

<div class="paragraph">

A job with two input parameters and one output can be defined like this:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
job, err := job.NewJob(
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

``` CodeRay
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

``` CodeRay
type ConcatString struct {
    A string
    B string
    S string
}

job, err := job.NewJob(
    WithKind("concat (one struct input and one struct output)"),
    WithExecutor(func(param *ConcatString) *ConcatString {
        // Store the concatenated string result in the input struct, and return it
        param.S = param.A + ", " + param.B
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

In this example, the result is also saved in the `S` field of the struct.

</div>

<div class="paragraph">

You can schedule jobs by passing arguments in various formats. For instance, to schedule a job using a struct as an argument:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
arg := &ConcatString{
    A: "Hello",
    B: "world!",
    S: "",
}
mgr.ScheduleJob(job, WithArguments(arg))
```

</div>

</div>

<div class="paragraph">

Alternatively, you can pass arguments as a JSON string:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
jsonArg := `{"A": "Hello", "B": "world!"}`
mgr.ScheduleJob(job, WithArguments(jsonArg))
```

</div>

</div>

<div class="paragraph">

`go-job` will automatically convert the provided arguments to the types expected by your job function. This means you can use the format that is most convenient for your application, whether it’s a struct, JSON, or other supported types.

</div>

</div>

<div class="sect3">

#### Executor with Auto-Injected Arguments

<div class="paragraph">

`go-job` supports special auto-injected arguments to allow you to access job context, manager, worker, and instance information directly within your function.

</div>

| Argument Type | Description | Example Usage in Executor Function |
|----|----|----|
| [context.Context](https://pkg.go.dev/context) | Context for cancellation and timeout | func(ctx context.Context) { …​ } |
| [job.Manager](https://pkg.go.dev/github.com/cybergarage/go-job/job#Manager) | Job manager interface | func(mgr job.Manager) { …​ } |
| [job.Worker](https://pkg.go.dev/github.com/cybergarage/go-job/job#Worker) | Worker processing the job instance | func(w job.Worker) { …​ } |
| [job.Instance](https://pkg.go.dev/github.com/cybergarage/go-job/job#Instance) | Current job instance | func(ji job.Instance) { …​ } |

<div class="paragraph">

These auto-injected arguments are automatically supplied when omitted. They enable you to control job execution, handle cancellation and timeouts, and access useful job-related metadata and methods.

</div>

<div class="sect4">

##### Using context.Context for Cancellation and Timeout

<div class="paragraph">

If your job function takes a `context.Context` argument, you can easily handle cancellation and timeout. For example, you can stop a long-running job when the context is canceled:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
job, err := job.NewJob(
    WithKind("sleep"),
    WithExecutor(func(ctx context.Context, d time.Duration) {
        steps := 10
        stepDuration := d / time.Duration(steps)
        for i := 0; i < steps; i++ {
            select {
            case <-ctx.Done():
                // Job is canceled or timed out
                return
            default:
                time.Sleep(stepDuration)
            }
        }
    }),
)
```

</div>

</div>

<div class="paragraph">

To schedule this job, pass the duration argument along with a dummy value (the actual context will be injected automatically):

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithArguments("?", time.Duration(1*time.Hour)))
```

</div>

</div>

<div class="paragraph">

Alternatively, you can specify only the duration argument because auto-injected arguments will be filled in automatically when omitted:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithArguments(time.Duration(1*time.Hour)))
```

</div>

</div>

</div>

<div class="sect4">

##### Accessing Manager, Worker, and Instance in Job Functions

<div class="paragraph">

You can also use auto-injected arguments like `job.Manager`, `job.Worker`, or `job.Instance` in your job function to get information about the job or control its execution.

</div>

<div class="paragraph">

For example, to log job details using the `Instance` argument:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
job, err := job.NewJob(
    WithKind("info (with special instance argument)"),
    WithExecutor(func(ji job.Instance) {
        // Log job details
        ji.Infof("%s (%s): attempts %d", ji.UUID(), ji.Kind(), ji.Attempts())
    }),
)
```

</div>

</div>

<div class="paragraph">

Schedule this job the same way:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithArguments("?"))
```

</div>

</div>

<div class="paragraph">

Or with the placeholder:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithArguments(job.Placeholder))
```

</div>

</div>

<div class="paragraph">

These features make it easy to write flexible and powerful job functions that can respond to cancellation, access job metadata, and interact with the job system.

</div>

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

``` CodeRay
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

``` CodeRay
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

``` CodeRay
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

``` CodeRay
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

``` CodeRay
job, err := job.NewJob(
    ....,
    WithCompleteProcessor(func(ji Instance, res []any) {
        ji.Infof("Result: %v", res)
    }),
    WithTerminateProcessor(func(ji Instance, err error) {
        ji.Errorf("Error: %v", err)
        return err
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

``` CodeRay
job, err := job.NewJob(
    ....,
    WithStateChangeProcessor(func(ji Instance, state JobState) {
        ji.Infof("State changed to: %v", state)
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

``` CodeRay
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

``` CodeRay
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

``` CodeRay
query := job.NewQuery(
    job.WithQueryInstance(ji), // filter by specific job instance
)
states := mgr.LookupInstanceHistory(query)
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

``` CodeRay
query := job.NewQuery(
    job.WithQueryInstance(ji), // filter by specific job instance
)
logs := mgr.LookupInstanceLogs(query)
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

<div class="sect3">

#### Setting Retry Policy

<div class="paragraph">

You can control how many times a job should be retried if it fails, allowing you to build more robust and fault-tolerant workflows. The retry policy can be set at both the job definition level and at the time of scheduling a job instance.

</div>

<div class="paragraph">

To set a default retry policy for all instances of a job, use `WithMaxRetries()` when creating the job. This determines the maximum number of times the job will be retried if it fails due to an error, timeout, or cancellation.

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
job, err := job.NewJob(
    WithMaxRetries(1), // Retry once if the job fails
)
```

</div>

</div>

<div class="paragraph">

You can also override the retry policy for a specific job instance at scheduling time by passing `WithMaxRetries()` to `ScheduleJob`. This is useful when you want to adjust the retry behavior for certain executions without changing the job’s default policy.

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithMaxRetries(2)) // Retry twice if the job fails
```

</div>

</div>

<div class="paragraph">

If a job instance fails, `go-job` will automatically retry it up to the specified number of times.

</div>

</div>

<div class="sect3">

#### Custom Termination Handling

<div class="paragraph">

You can define custom logic to handle job termination using `WithTerminateProcessor()`. This allows you to inspect the error, decide whether to retry, transform the error, or perform cleanup actions.

</div>

<div class="paragraph">

For example, you might want to avoid retrying only for specific errors (such as a timeout):

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job,
    WithTerminateProcessor(func(ji Instance, err error) error {
        if errors.Is(err, context.DeadlineExceeded) {
            // Do not retry if the job was terminated due to a deadline being exceeded
            ji.Infof("Job (%s) terminated due to deadline exceeded: %v", ji.Kind(), err)
            return nil
        }
        // Retry for all other errors
        return error
    }),
)
```

</div>

</div>

<div class="paragraph">

The `WithTerminateProcessor` function is called whenever a job instance ends with an error. You can return `nil` to mark the job as successful, return the original error to keep it as failed, or return a new error to transform the result.

</div>

</div>

<div class="sect3">

#### Setting Backoff Strategy

<div class="paragraph">

A backoff strategy controls how long the system waits before retrying a failed job. This helps prevent overwhelming the system or external resources with rapid retries.

</div>

<div class="paragraph">

You can set a custom backoff strategy using `WithBackoffStrategy()`. The function you provide should return the duration to wait before the next retry attempt.

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
mgr.ScheduleJob(job, WithBackoffStrategy(
    func(ji Instance) time.Duration {
        // Exponential backoff
                return time.Duration(float64(ji.Attempts()) * float64(time.Second) * (0.8 + 0.4*rand.Float64()))
    },
))
```

</div>

</div>

<div class="paragraph">

You can implement more advanced strategies, such as exponential backoff, by adjusting the returned duration based on the number of attempts or other factors.

</div>

<div class="paragraph">

These features allow you to build robust, fault-tolerant job processing pipelines that can gracefully handle

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

``` CodeRay
// High priority job (executed first)
highPriorityJob, err := job.NewJob(
    WithKind("urgent-task"),
    WithPriority(0), // lower number = higher priority like Unix nice values
    WithExecutor(func() { fmt.Println("Urgent task executing") }),
)

// Low priority job (executed later)
lowPriorityJob, err := job.NewJob(
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

``` CodeRay
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

``` CodeRay
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

``` CodeRay
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

``` CodeRay
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

### System Jobs

<div class="paragraph">

`go-job` comes with several built-in system jobs that are ready to use. You can schedule and configure these jobs just like any other job in your application.

</div>

<div class="paragraph">

The following table lists each system job and its purpose:

</div>

| Function | Kind | Description |
|----|----|----|
| [NewHistoryCleaner](../job/plugins/job/system/history_cleaner.go) | system.history.cleaner | Deletes old job history records |
| [NewLogCleaner](../job/plugins/job/system/log_cleaner.go) | system.log.cleaner | Deletes old job log records |

<div class="paragraph">

For step-by-step usage and scheduling examples, see [Extension Guide](extension-guide.md) documentation or [the usage examples](../jobtest/example_job_plugins_test.go).

</div>

</div>

<div class="sect2">

### Remote Management with gRPC API

<div class="paragraph">

`go-job` provides a comprehensive gRPC API for remote job management, enabling you to schedule, monitor, and control jobs from external systems or distributed environments. This allows seamless integration with microservices, orchestration platforms, and remote applications.

</div>

<div class="imageblock">

<div class="content">

![job framework](img/job-framework.png)

</div>

</div>

<div class="sect3">

#### Server vs. Manager: Which to Use?

<div class="paragraph">

`go-job` provides two main components:

</div>

<div class="ulist">

- **Manager**: Embed this in your Go application to schedule and run jobs directly in your process. Use Manager when you want to control jobs from your own code, such as in microservices, batch jobs, or automation tools.

- **Server**: Run as a standalone process to provide job management over a gRPC API. Use Server when you want to manage jobs remotely, integrate with other systems, or offer job scheduling as a service for multiple applications.

</div>

<div class="paragraph">

When to use each:

</div>

<div class="ulist">

- **Manager**: For direct, in-process job control and tight integration with your Go code.

- **Server**: For remote job management, gRPC access, or centralized job scheduling across many apps.

</div>

<div class="paragraph">

The Server uses a Manager internally, so all features are available.

</div>

</div>

<div class="sect3">

#### Remote Operation with gRPC API

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

`go-job` makes it easy to build distributed systems by allowing you to plug in different storage backends using the `Store` interface. This means that multiple `go-job` instances can work together and share job data through a common store, as shown in the diagram below:

</div>

<div class="imageblock">

<div class="content">

![job store](img/job-store.png)

</div>

</div>

<div class="paragraph">

By choosing or implementing a suitable store (such as etcd or Valkey), you can ensure that job metadata and execution state are accessible from any node in your system.

</div>

<div class="paragraph">

This approach enables:

</div>

<div class="ulist">

- Distributed job scheduling across multiple nodes

- Coordination of jobs between different servers

- Persistence of job state even if a node restarts

- Fault-tolerant execution, so jobs are not lost if a node fails

</div>

<div class="paragraph">

`go-job` comes with several built-in store plugins you can use right away:

</div>

| Store | Version | Driver | Type | Persistence | Distribution | Use Case | Notes |
|----|----|----|----|----|----|----|----|
| [Valkey](https://valkey.io/) | \>=8.1.3 | [valkey-go](https://github.com/valkey-io/valkey-go) v1.0.63 | External (Valkey) | Optional | Yes | Production/Distributed | Redis-compatible |
| [Redis](https://redis.io/) | \>=7.2.4 | [go-redis](https://github.com/redis/go-redis/) v9.12.1 | External (Redis) | Optional | Yes | Production/Distributed | Popular in-memory store |
| [etcd](https://etcd.io/) | \>=3.6.4 | [etcd/client](https://pkg.go.dev/go.etcd.io/etcd/client/v3) v3.0.1 | External (etcd) | Yes | Yes | Production/Distributed | Strong consistency |
| [go-memdb](https://github.com/hashicorp/go-memdb/) | \>=1.3.5 | (none) | In-memory | No | No | Testing/Development | Fastest but data is lost on restart |

<div class="paragraph">

For more details about the `Store` interface and how to extend it, see the [Design and Architecture](design.md) and [Extension Guide](extension-guide.md) documentation.

</div>

<div class="sect4">

##### Valkey Store Plugin

<div class="paragraph">

`Valkey` is a fast and lightweight key-value store built on the Valkey library. It offers a simple and efficient way to store and retrieve job metadata and state in `go-job`.

</div>

<div class="paragraph">

To use the Valkey store plugin, create a manager instance with Valkey as the backend:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
import (
    "net"

    "github.com/cybergarage/go-job/job"
    "github.com/cybergarage/go-job/job/plugins/store"
    "github.com/valkey-io/valkey-go"
)

func main() {
    valkeyOpt := valkey.ClientOption{
        InitAddress: []string{net.JoinHostPort("10.0.0.10", "6379")},
    }
    mgr, err := job.NewManager(
        job.WithStore(store.NewValkeyStore(valkeyOpt)),
    )
}
```

</div>

</div>

</div>

<div class="sect4">

##### Etcd Store Plugin

<div class="paragraph">

`etcd` is a distributed key-value store used to manage job metadata and state in `go-job`. It is built on the etcd v3 API and provides advanced features such as lease-based locking and real-time notifications using the watch mechanism.

</div>

<div class="paragraph">

To use the etcd store plugin, simply create a new manager instance with etcd as the backend:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
import (
    "net"

    "github.com/cybergarage/go-job/job"
    "github.com/cybergarage/go-job/job/plugins/store"
    v3 "go.etcd.io/etcd/client/v3"
)

func main() {
        etcdOpt := v3.Config{
                Endpoints: []string{net.JoinHostPort("10.0.0.10", "2379")},
        }
        mgr, err := job.NewManager(
                job.WithStore(store.NewEtcdStore(etcdOpt)),
        )
}
```

</div>

</div>

</div>

</div>

</div>

</div>

</div>

<div id="footer">

<div id="footer-text">

Last updated 2025-08-27 23:35:07 +0900

</div>

</div>
