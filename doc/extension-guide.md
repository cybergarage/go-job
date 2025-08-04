# Extension Guide 
---
lang: en
title: Extension Guide
viewport: width=device-width, initial-scale=1.0
---

<div id="header">

# Extension Guide

</div>

<div id="content">

<div id="preamble">

<div class="sectionbody">

<div class="paragraph">

This guide provides an overview of how to extend `go-job` with custom plugins, allowing you to add new functionality or integrate with external systems.

</div>

</div>

</div>

<div class="sect1">

## Plugin Development

<div class="sectionbody">

<div class="sect2">

### Store Plugin

<div class="sect3">

#### Store Interface

<div class="imageblock">

<div class="content">

![job store](img/job-store.png)

</div>

</div>

<div class="listingblock">

<div class="content">

``` highlight
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
    // DequeueInstance removes a job instance from the store by its unique identifier.
    DequeueInstance(ctx context.Context, job Instance) error
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
    // LookupInstanceHistory lists all state records for a job instance. The returned history is sorted by their timestamp.
    LookupInstanceHistory(ctx context.Context, job Instance) (InstanceHistory, error)
    // ListInstanceHistory lists all state records for all job instances. The returned history is sorted by their timestamp.
    ListInstanceHistory(ctx context.Context) (InstanceHistory, error)
    // ClearInstanceHistory clears all state records for a job instance.
    ClearInstanceHistory(ctx context.Context) error
}

// LogStore is an interface that defines methods for logging job instance messages.
type LogStore interface {
    // Infof logs an informational message for a job instance.
    Infof(ctx context.Context, job Instance, format string, args ...any) error
    // Warnf logs a warning message for a job instance.
    Warnf(ctx context.Context, job Instance, format string, args ...any) error
    // Errorf logs an error message for a job instance.
    Errorf(ctx context.Context, job Instance, format string, args ...any) error
    // LookupInstanceLogs lists all log entries for a job instance. The returned logs are sorted by their timestamp.
    LookupInstanceLogs(ctx context.Context, job Instance) ([]Log, error)
    // ClearInstanceLogs clears all log entries for a job instance.
    ClearInstanceLogs(ctx context.Context) error
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

Last updated 2025-08-04 22:16:03 +0900

</div>

</div>
