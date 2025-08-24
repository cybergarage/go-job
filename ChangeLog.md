# ChangeLog

## 1.2.1 (2025-0X-XX)
### ✨ New Features
- **System Worker**
  - Store cleaner plugin
### 🛠 Enhancements
- **Query**
  - Limit and offset support

## 1.2.0 (2025-08-24)
## 🚀 Features
- **Job Cancellation & Control**
  - Enhanced `jobctl` with new `cancel` and `schedule` commands for better job management.
    ([b8f7ca8](https://github.com/cybergarage/go-job/commit/b8f7ca8), [1d2e16d](https://github.com/cybergarage/go-job/commit/1d2e16d), [89f08e5](https://github.com/cybergarage/go-job/commit/89f08e5), [83ca91a](https://github.com/cybergarage/go-job/commit/83ca91a))
  - Added `JobService::CancelInstances()` and related CLI support to cancel specific jobs or instances.
    ([7bd8377](https://github.com/cybergarage/go-job/commit/7bd8377), [c3c1bc](https://github.com/cybergarage/go-job/commit/3c1c1bc))
  - Added `Manager::CancelInstances()` to cancel all job instances matching a query.
    ([3c1c1bc](https://github.com/cybergarage/go-job/commit/3c1c1bc), [d9b33f2](https://github.com/cybergarage/go-job/commit/d9b33f2))
  - Improved job state tracking for cancellation and timeout.
    ([42a2779](https://github.com/cybergarage/go-job/commit/42a2779), [2551496](https://github.com/cybergarage/go-job/commit/2551496), [e94e067](https://github.com/cybergarage/go-job/commit/e94e067))
- **Worker & Manager Enhancements**
  - Added `WorkerGroup::Wait()`, `Worker::Wait()`, and `Manager::Wait()` to synchronize job completion and termination.
    ([e113453](https://github.com/cybergarage/go-job/commit/e113453), [72662b8](https://github.com/cybergarage/go-job/commit/72662b8), [d54ab87](https://github.com/cybergarage/go-job/commit/d54ab87))
  - Added context support for timeout/cancel to `Wait` and `ResizeWorkers` methods.
    ([896c3bc](https://github.com/cybergarage/go-job/commit/896c3bc))
  - Added various methods for clearing instance logs and history.
    ([7b8d216](https://github.com/cybergarage/go-job/commit/7b8d216), [22937ca](https://github.com/cybergarage/go-job/commit/22937ca))
- **Configuration & Client**
  - Added a new `Config` interface and integrated config handling for `jobd`.
    ([fea290a](https://github.com/cybergarage/go-job/commit/fea290a), [7e7a762](https://github.com/cybergarage/go-job/commit/7e7a762))
  - Improved `CLIClient.Execute` error output with detailed messages.
    ([b2eef37](https://github.com/cybergarage/go-job/commit/b2eef37))
- **Prometheus Metrics**
  - Added Prometheus metrics server to monitor job lifecycle and worker status.  
    ([3cd693e](https://github.com/cybergarage/go-job/commit/3cd693e), [688f295](https://github.com/cybergarage/go-job/commit/688f295))
## 🛠 Fixes and Improvements
- Improved cancellation and timeout handling in `Execute` and related methods.
  ([9901727](https://github.com/cybergarage/go-job/commit/9901727), [0f78ee5](https://github.com/cybergarage/go-job/commit/0f78ee5), [66872e9](https://github.com/cybergarage/go-job/commit/66872e9))
- Fixed worker state handling for cancellations and timeouts.
  ([f5e2950](https://github.com/cybergarage/go-job/commit/f5e2950), [a67b785](https://github.com/cybergarage/go-job/commit/a67b785), [39e63ae](https://github.com/cybergarage/go-job/commit/39e63ae), [690a9ca](https://github.com/cybergarage/go-job/commit/690a9ca))
- Fixed etcd client handling in plugin store startup.
  ([3b16dac](https://github.com/cybergarage/go-job/commit/3b16dac))
- Refactored and renamed variables for consistency (e.g., `JobCancelled` to `JobCanceled`).
  ([e94e067](https://github.com/cybergarage/go-job/commit/e94e067), [67b9189](https://github.com/cybergarage/go-job/commit/67b9189))

## 1.1.1 (2025-08-20)
### 🛠 Enhancements
- **Job Management**
  - Executor supports special instance arguments
- **Logging**
  - Added debug log level and Debugf methods

## 1.1.0 (2025-08-18)
### ✨ New Features
- **Job Management**
  - Job can handle policy options
  - `Process()` now supports timeout policy option
### 🛠 Enhancements
- **Query**
  - Added `WithQueryBefore()` and `WithQueryAfter()` for time-based filtering
- **Store Interface**
  - Added Redis plugin implemented for distributed job management
### 🐛 Bug Fixes
- **Manager**
  - `WithJob()` now checks if the passed job is nil
  - `Registry::LookupJob()` checks whether looked-up job is nil

## 1.0.0 (2025-08-13)
### ✨ Features
- **gRPC API**: Added full-featured gRPC service enabling remote job operations
  - Job scheduling and execution management
  - Real-time job monitoring and status tracking
  - Cross-platform client support
- **CLI Tool**: Introduced `jobctl` for command-line job administration
  - Schedule jobs remotely with arguments and timing options
  - Query job instances by kind, state, or time range
  - List and monitor registered jobs
- **Store Interface**: Introduced Store interface for better abstraction and flexibility in distributed job management
  - **kv.Store interface**: Introduced a new `kv.Store` interface for key-value storage operations, allowing for flexible backend implementations.
    - Added etcd and Valkey plugins for distributed job management.
    - Added memdb plugins for testing.

## v0.9.3 (2025-07-26)
### 🛠 Enhancements
- **Store Interface**: Introduced pluggable Store interface enabling flexible distributed job management across different storage backends
- Updated scheduling logic to automatically schedule jobs with timing configuration upon registration.
- Improved QueueStore interface with `DequeueNextInstance` for priority-based job retrieval.

## v0.9.2 (2025-07-25)
### 🛠 Enhancements
- **Instance Management**: Added `Manager::ListInstances()`, `LookupInstances()`, and `Query` system for comprehensive instance filtering and retrieval
- **Data Structures**: Enhanced `InstanceMap` with `Arguments()`, `ResultSet()`, `Error()` methods and improved argument handling with JSON support
- **State Management**: Added `JobStateAll` constant, renamed `JobStateUnknown` to `JobStateUnset`, improved timestamp tracking
- **API Improvements**: Extended Job/Instance interfaces with metadata support (`Description()`, `RegisteredAt()`, `ProcessedAt()`) and better lifecycle management
- **Worker & Queue**: Enhanced error handling during termination, improved rescheduling, added `Queue::List()` method
- **Protocol Updates**: Multiple service.proto enhancements for better client-server communication

## v0.9.1 (2025-07-20)
### ✨ Features
- Extended `Manager` to clear instance logs and history
- Extended `Instance` interface with `Logs()` method
### 🛠 Enhancements
- Extended `Log` interface with `Kind()` method
- Enhanced `Worker::Run()` to log job errors
- Enabled backoff strategy in `Worker::Run()`

## v0.9.0 (2025-07-17)
- Initial release
### Major Features
- Flexible job creation and registration with custom executors and arguments
- Priority and policy management for job execution order
- Job execution with response and error handlers
- Job state and history logging for monitoring and debugging
