<div id="header">

# Extension Guide

</div>

<div id="content">

<div id="preamble">

<div class="sectionbody">

<div class="paragraph">

This guide provides an overview of how to extend `go-job` with custom plugins, allowing you to add new functionality or integrate with external systems.

</div>

<div id="toc" class="toc">

<div id="toctitle" class="title">

Table of Contents:

</div>

- [Job Plugin Development](#_job_plugin_development)
  - [Job Interface](#_job_interface)
  - [Special Arguments for Executor](#_special_arguments_for_executor)
  - [Built-in System Jobs](#_built_in_system_jobs)
- [Store Plugin Development](#_store_plugin_development)
  - [Store Interface](#_store_interface)
  - [kv.Store Interface](#_kv_store_interface)
    - [Valkey Store Plugin](#_valkey_store_plugin)
    - [Etcd Store Plugin](#_etcd_store_plugin)

</div>

</div>

</div>

<div class="sect1">

## Job Plugin Development

<div class="sectionbody">

<div class="paragraph">

`go-job` provides several built-in system jobs that help you keep your job server clean and running smoothly—no extra setup required. You can also extend `go-job` by creating your own job plugins.

</div>

<div class="sect3">

#### Job Interface

<div class="paragraph">

To create a custom job plugin, simply implement the following interface:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
package job

// Job represents a job that can be scheduled to run at a specific time or interval.
type Job interface {
        // Kind returns the name of the job.
        Kind() string
        // Description returns a description of the job.
        Description() string
        // Handler returns the job handler for the job.
        Handler() Handler
        // Schedule returns the schedule for the job.
        Schedule() Schedule
        // Policy returns the policy for the job.
        Policy() Policy
        // RegisteredAt returns the time when the job was registered.
        RegisteredAt() time.Time
        // Map returns a map representation of the job.
        Map() map[string]any
        // String returns a string representation of the job.
        String() string
}
```

</div>

</div>

</div>

<div class="sect3">

#### Special Arguments for Executor

<div class="paragraph">

Executor functions can use special arguments to easily access job context, manager, worker, and instance information. This makes it simple to control job execution, handle cancellation and timeout, and access useful metadata.

</div>

| Argument Type | Description | Example Usage in Executor Function |
|----|----|----|
| [context.Context](https://pkg.go.dev/context) | Context for cancellation and timeout | func(ctx context.Context) { …​ } |
| [job.Manager](https://pkg.go.dev/github.com/cybergarage/go-job/job#Manager) | Job manager interface | func(mgr job.Manager) { …​ } |
| [job.Worker](https://pkg.go.dev/github.com/cybergarage/go-job/job#Worker) | Worker processing the job instance | func(w job.Worker) { …​ } |
| [job.Instance](https://pkg.go.dev/github.com/cybergarage/go-job/job#Instance) | Current job instance | func(ji job.Instance) { …​ } |

<div class="paragraph">

These special arguments allow your job logic to be more flexible and powerful.

</div>

</div>

<div class="sect3">

#### Built-in System Jobs

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

For step-by-step usage and scheduling examples, see [the usage examples](../jobtest/example_job_plugins_test.go).

</div>

</div>

</div>

</div>

<div class="sect1">

## Store Plugin Development

<div class="sectionbody">

<div class="paragraph">

The `go-job` framework supports custom store plugins that can be used to manage job instances, their states, and logs. A store plugin must implement the `Store` interface, which defines methods for managing job instances and their histories.

</div>

<div class="imageblock">

<div class="content">

![job store](img/job-store.png)

</div>

</div>

<div class="sect2">

### Store Interface

<div class="paragraph">

The Store interface specifies the required methods that every store plugin must implement to manage job instances, their states, and logs:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
package job

// Store defines the interface for job queue, history, and logging.
type Store interface {
    // Name returns the name of the store.
    Name() string
    // PendingStore provides methods for managing job instances.
    QueueStore
    // HistoryStore provides methods for managing job instance state history.
    HistoryStore
    // Start starts the store.
    Start() error
    // Stop stops the store.
    Stop() error
}

// QueueStore is an interface that defines methods for managing job instances in a pending state.
type QueueStore interface {
    // EnqueueInstance stores a job instance in the store.
    EnqueueInstance(ctx context.Context, job Instance) error
    // DequeueNextInstance retrieves and removes the highest priority job instance from the store. If no job instance is available, it returns nil.
    DequeueNextInstance(ctx context.Context) (Instance, error)
    // ListInstances lists all job instances in the store.
    ListInstances(ctx context.Context) ([]Instance, error)
    // ClearInstances clears all job instances in the store.
    ClearInstances(ctx context.Context) error
}

// HistoryStore is an interface that defines methods for managing job instance state history.
type HistoryStore interface {
    // StateStore provides methods for managing job instance state history.
    StateStore
    // LogStore provides methods for logging job instance messages.
    LogStore
}

// StateStore is an interface that defines methods for managing job instance state history.
type StateStore interface {
    // LogInstanceState adds a new state record for a job instance.
    LogInstanceState(ctx context.Context, state InstanceState) error
    // LookupInstanceHistory lists all state records for a job instance that match the specified query. The returned history is sorted by their timestamp.
    LookupInstanceHistory(ctx context.Context, query Query) (InstanceHistory, error)
    // ClearInstanceHistory clears all state records for a job instance that match the specified filter.
    ClearInstanceHistory(ctx context.Context, filter Filter) error
}

// LogStore is an interface that defines methods for logging job instance messages.
type LogStore interface {
    // Infof logs an informational message for a job instance.
    Infof(ctx context.Context, job Instance, format string, args ...any) error
    // Warnf logs a warning message for a job instance.
    Warnf(ctx context.Context, job Instance, format string, args ...any) error
    // Errorf logs an error message for a job instance.
    Errorf(ctx context.Context, job Instance, format string, args ...any) error
    // LookupInstanceLogs lists all log entries for a job instance that match the specified query. The returned logs are sorted by their timestamp.
    LookupInstanceLogs(ctx context.Context, query Query) ([]Log, error)
    // ClearInstanceLogs clears all log entries for a job instance that match the specified filter.
    ClearInstanceLogs(ctx context.Context, filter Filter) error
}
```

</div>

</div>

</div>

<div class="sect2">

### kv.Store Interface

<div class="paragraph">

To create a custom store plugin using a key-value store, `go-job` provides a straightforward key-value store interface.

</div>

<div class="paragraph">

This interface makes it easy to build your own plugins for storing and managing job data.

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
// Store represents a key-value store interface.
type Store interface {
    // UniqueKeys returns whether keys should be unique.
    UniqueKeys() bool
    // Name returns the name of the store.
    Name() string
    // Set stores a key-value object. If the key already holds some value, it is overwritten.
    Set(ctx context.Context, obj Object) error
    // Get returns a key-value object of the specified key.
    Get(ctx context.Context, key Key) (Object, error)
    // Scan returns a result set of all key-value objects whose keys have the specified prefix.
    Scan(ctx context.Context, key Key, opts ...Option) (ResultSet, error)
    // Remove removes the specified key-value object.
    Remove(ctx context.Context, obj Object) error
    // Delete deletes all key-value objects whose keys have the specified prefix.
    Delete(ctx context.Context, key Key) error
    // Dump returns all key-value objects in the store.
    Dump(ctx context.Context) ([]Object, error)
    // Start starts the store.
    Start() error
    // Stop stops the store.
    Stop() error
    // Clear removes all key-value objects from the store.
    Clear() error
}
```

</div>

</div>

<div class="paragraph">

By default, `go-job` provides ready-to-use key-value store implementations, including Valkey and Etcd.

</div>

<div class="paragraph">

The following table summarizes the main differences between the available store plugins:

</div>

| Store | Version | Driver | Type | Persistence | Distribution | Use Case | Notes |
|----|----|----|----|----|----|----|----|
| [Valkey](https://valkey.io/) | \>=8.1.3 | [valkey-go](https://github.com/valkey-io/valkey-go) v1.0.63 | External (Valkey) | Optional | Yes | Production/Distributed | Redis-compatible |
| [Redis](https://redis.io/) | \>=7.2.4 | [go-redis](https://github.com/redis/go-redis/) v9.12.1 | External (Redis) | Optional | Yes | Production/Distributed | Popular in-memory store |
| [etcd](https://etcd.io/) | \>=3.6.4 | [etcd/client](https://pkg.go.dev/go.etcd.io/etcd/client/v3) v3.0.1 | External (etcd) | Yes | Yes | Production/Distributed | Strong consistency |
| [go-memdb](https://github.com/hashicorp/go-memdb/) | \>=1.3.5 | (none) | In-memory | No | No | Testing/Development | Fastest but data is lost on restart |

<div class="paragraph">

These store plugins are built on top of the `kv.Store` interface, making it easy to use them out of the box or to develop your own custom store plugins based on the same interface.

</div>

<div class="paragraph">

To use one of these built-in stores, simply create a manager instance and specify the desired backend (Valkey or Etcd) when configuring the store:

</div>

<div class="sect4">

##### Valkey Store Plugin

<div class="paragraph">

To use the Valkey store plugin, simply create a manager instance with Valkey as the backend:

</div>

<div class="listingblock">

<div class="content">

``` CodeRay
import (
    "net"

    "github.com/cybergarage/go-job/job"
    "github.com/cybergarage/go-job/job/plugins/store"
    "github.com/cybergarage/go-job/job/plugins/store/kv/valkey"
    v1 "github.com/valkey-io/valkey-go"
)

func main() {
    valkeyOpt := v1.ClientOption{
        InitAddress: []string{net.JoinHostPort("10.0.0.10", "6379")},
    }
    mgr, err := job.NewManager(
        job.WithStore(store.NewKvStoreWith(valkey.NewStore(valkeyOpt))),
    )
}
```

</div>

</div>

</div>

<div class="sect4">

##### Etcd Store Plugin

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
    "github.com/cybergarage/go-job/job/plugins/store/kv/etcd"
    v3 "go.etcd.io/etcd/client/v3"
)

func main() {
    etcdOpt := v3.Config{
        Endpoints: []string{net.JoinHostPort("10.0.0.10", "2379")},
    }
    mgr, err := job.NewManager(
        job.WithStore(store.NewKvStoreWith(etcd.NewStore(etcdOpt))),
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

Last updated 2025-08-26 13:23:32 +0900

</div>

</div>
