<div id="header">

# Design and Architecture

</div>

<div id="content">

<div id="preamble">

<div class="sectionbody">

<div class="paragraph">

This document provides a detailed overview of \`go-job’s design and architecture, including future plans.

</div>

<div id="toc" class="toc">

<div id="toctitle" class="title">

Table of Contents:

</div>

- [Design Concept](#_design_concept)
- [Terminology](#_terminology)
- [Key Components](#_key_components)
- [Deployment Options](#_deployment_options)
- [Job State Lifecycle](#_job_state_lifecycle)
- [Job Registration and Processing Flow](#_job_registration_and_processing_flow)
- [Prometheus Metrics](#_prometheus_metrics)

</div>

</div>

</div>

<div class="sect1">

## Design Concept

<div class="sectionbody">

<div class="paragraph">

These design choices position `go-job` as a scalable, flexible, and production-grade job management framework for Go applications that require robust orchestration of asynchronous or scheduled tasks.

</div>

<div class="paragraph">

The core design principles of `go-job` are:

</div>

<div class="ulist">

- **General-Purpose Job Definition**: `go-job` allows registration of arbitrary Go functions with any signature using `any`-typed arguments and return values.

- **Flexible Scheduling**: `go-job` supports not only cron-style and fixed-interval execution, but also delayed, timed, and immediate scheduling — all within a unified API.

- **Distributed and Observable by Design**: `go-job` introduces a pluggable `Store` interface to enable consistent state sharing and coordination between nodes, while also providing first-class support for logging, state transitions, and job lifecycle monitoring.

- **Extensibility First**: Every component — executors, stores, workers, handlers — is designed to be pluggable or replaceable, making `go-job` suitable for embedded use, microservices, and server-mode deployment with gRPC APIs.

</div>

<div class="paragraph">

For a comparison of design concepts with other OSS job frameworks, see [go-job Comparison (OpenAI Research)](https://github.com/cybergarage/go-job/blob/main/doc/design-comparison.md).

</div>

</div>

</div>

<div class="sect1">

## Terminology

<div class="sectionbody">

<div class="paragraph">

This section defines the key terms and concepts used throughout the `go-job` system. Understanding these terms is essential for working effectively with the job scheduling and execution framework.

</div>

| Term | Definition |
|----|----|
| Job | A reusable definition that specifies work to be performed, including the executor function, scheduling rules, and retry policies. |
| Job Instance | A specific execution of a job with concrete arguments, unique identifier, and state tracking throughout its lifecycle. |
| Executor | A Go function that implements the actual business logic for a specific job type. |
| Processor | A function that processes job instances, including completion and termination logic. |
| Worker | A component that executes job instances by invoking the registered executors, handling retries, and managing state transitions. |

</div>

</div>

<div class="sect1">

## Key Components

<div class="sectionbody">

<div class="paragraph">

`go-job` is designed to handle job scheduling and execution efficiently. The architecture consists of several key components that work together to provide a robust job processing system.

</div>

<div class="paragraph">

The main components of `go-job` are:

</div>

<div class="imageblock">

<div class="content">

![job framework](img/job-framework.png)

</div>

</div>

| Component | Description |
|----|----|
| Server | Provides gRPC endpoints for job scheduling and management, allowing clients to interact with the job system through the manager. |
| Manager | Coordinates job scheduling and execution across go-job components. |
| Registry | Holds job definitions and their associated executors. |
| Worker | Processes job instances by executing the registered functions. |
| Queue | Manages job instances, ensuring they are processed in the correct order. |
| History | Tracks state transitions of job instances, providing an execution history. |
| Log | Captures logs for each job instance, providing detailed execution information. |
| Store | Provides abstracted persistence for job metadata and execution state, enabling distributed operation and fault tolerance. |

</div>

</div>

<div class="sect1">

## Deployment Options

<div class="sectionbody">

<div class="paragraph">

go-job can be used in two main ways: by embedding the manager directly in your Go application, or by running the standalone server component. This section explains the differences between these deployment options and helps you choose the best approach for your use case.

</div>

| Feature | Manager | Server |
|----|----|----|
| Usage | Embedded in your Go application | Standalone process accessible over the network |
| Job Scheduling | Yes (within the application) | Yes (via gRPC API from clients) |
| Job State & Log Management | Yes (within the application) | Yes (remotely via gRPC API) |
| Metrics Export | No | Yes (Prometheus-compatible endpoint) |
| Best For | Standalone/local job management | Distributed/remote job management and monitoring |
| API Exposure | No external API | gRPC API and Prometheus metrics endpoint |

<div class="sect2">

### Selecting Manager Usage

<div class="paragraph">

To use go-job, you can embed the manager directly in your Go application to schedule jobs, manage job instances, and process their states and logs. This approach allows you to handle all job management tasks easily within your application.

</div>

<div class="paragraph">

For more information about embedding the manager in your Go application, see the [Quick Start](quick-start.md) and [Go Reference](https://pkg.go.dev/github.com/cybergarage/go-job) documentation.

</div>

</div>

<div class="sect2">

### Selecting Server Usage

<div class="paragraph">

Alternatively, you can use the go-job server component, which provides a gRPC interface for remote job management. This enables clients to schedule jobs and retrieve job states and logs over the network.

</div>

<div class="paragraph">

In addition to the gRPC API, the server also exposes a Prometheus-compatible metrics endpoint. This allows you to monitor job execution, server health, and other key metrics using standard Prometheus tools.

</div>

<div class="paragraph">

For more information about the server component, see the [gRPC API](grpc-api.md) and [CLI (jobctl)](cmd/cli/jobctl.md) documentation.

</div>

</div>

</div>

</div>

<div class="sect1">

## Job State Lifecycle

<div class="sectionbody">

<div class="paragraph">

The job state in `go-job` is managed through a combination of job instances and their associated states. The state of a job instance is crucial for understanding its lifecycle and for debugging purposes.

</div>

<div class="imageblock">

<div class="content">

![job state](img/job-state.png)

</div>

</div>

| Category | State | Description |
|----|----|----|
| Initial | Created | The job instance has been created and is awaiting scheduling. |
| Active | Scheduled | The job instance has been queued and is waiting to be processed by a worker. |
| Active | Processing | The job instance is currently being executed by a worker. |
| Final | Terminated | The job instance encountered an error or was forcibly stopped before completion. |
| Final | Completed | The job instance finished successfully. |
| Final | Canceled | The job instance was canceled before completion. |
| Final | TimedOut | The job instance timed out before completion. |

<div class="paragraph">

Each job instance can transition through various states, such as `Scheduled`, `Processing`, `Completed`, `Terminated`, `Canceled`, and `TimedOut`. These states are tracked in the job manager, allowing you to monitor the progress and outcome of each job instance.

</div>

</div>

</div>

<div class="sect1">

## Job Registration and Processing Flow

<div class="sectionbody">

<div class="paragraph">

The `go-job` server is designed to be modular and extensible. Each component, including the registry, manager, and worker, can be independently developed and maintained.

</div>

<div class="paragraph">

The following sequence diagram illustrates the flow of job registration, scheduling, and processing.

</div>

<div class="imageblock">

<div class="content">

![job seqdgm](img/job-seqdgm.png)

</div>

</div>

<div class="sect2">

### Store Plugins and Registry Sharing Limitations

<div class="paragraph">

Currently, the registry that holds job definitions cannot be shared between go-job servers. Because Go does not support serializing or transmitting function pointers (executors) over RPC, each go-job server must maintain its own local registry of job definitions.

</div>

<div class="imageblock">

<div class="content">

![job store](img/job-store.png)

</div>

</div>

<div class="quoteblock">

> <div class="paragraph">
>
> **Note:** In the future, support for sharing the registry across go-job servers may be added through technologies such as shell scripts, Python, and WebAssembly (Wasm), but there are currently no concrete plans for this feature.
>
> </div>

</div>

<div class="paragraph">

The queue, history, and log components can be shared between go-job servers using distributed store plugins. The following table summarizes the main differences between the available store plugins:

</div>

| Store | Version | Driver | Type | Persistence | Distribution | Use Case | Notes |
|----|----|----|----|----|----|----|----|
| [Valkey](https://valkey.io/) | \>=8.1.3 | [valkey-go](https://github.com/valkey-io/valkey-go) v1.0.63 | External (Valkey) | Optional | Yes | Production/Distributed | Redis-compatible |
| [Redis](https://redis.io/) | \>=7.2.4 | [go-redis](https://github.com/redis/go-redis/) v9.12.1 | External (Redis) | Optional | Yes | Production/Distributed | Popular in-memory store |
| [etcd](https://etcd.io/) | \>=3.6.4 | [etcd/client](https://pkg.go.dev/go.etcd.io/etcd/client/v3) v3.0.1 | External (etcd) | Yes | Yes | Production/Distributed | Strong consistency |
| [go-memdb](https://github.com/hashicorp/go-memdb/) | \>=1.3.5 | (none) | In-memory | No | No | Testing/Development | Fastest but data is lost on restart |

<div class="paragraph">

This enables a distributed architecture where multiple go-job servers can operate together, sharing job instances and state information. To learn more about the store plugins, see [Extension Guide](extension-guide.md).

</div>

</div>

</div>

</div>

<div class="sect1">

## Prometheus Metrics

<div class="sectionbody">

<div class="paragraph">

`go-job` includes a built-in Prometheus metrics endpoint that exposes a variety of metrics related to the job server’s performance and health.

</div>

<div class="imageblock">

<div class="content">

![job framework](img/job-framework.png)

</div>

</div>

<div class="sect2">

### Observability and Metrics

| Metric Name | Type | Labels | Description |
|----|----|----|----|
| go_job_registered | Gauge |  | Current number of registered jobs |
| go_job_workers | Gauge |  | Current number of workers |
| go_job_queued | GaugeVec | kind | Current number of queued jobs by kind |
| go_job_executed_total | CounterVec | kind | Total number of executed jobs by kind |
| go_job_completed_total | CounterVec | kind | Total number of successfully completed jobs by kind |
| go_job_terminated_total | CounterVec | kind | Total number of terminated jobs by kind |
| go_job_canceled_total | CounterVec | kind | Total number of canceled jobs by kind |
| go_job_timedout_total | CounterVec | kind | Total number of timed out jobs by kind |
| go_job_duration_seconds | Histogram | kind | Histogram of job execution durations in seconds by kind |

</div>

</div>

</div>

</div>

<div id="footer">

<div id="footer-text">

Last updated 2025-08-27 08:21:58 +0900

</div>

</div>
