# ChangeLog

## 1.2.0 (2025-0X-XX)
### ✨ New Features
- **Job Management**
  - Job instance cancellation
- **System Worker**
  - Store cleaner support
- **Metrics**
  - Prometheus metrics support added
### 🛠 Enhancements
- **Query**
  - Limit and offset support

## 1.1.1 (2025-0X-XX)
### 🛠 Enhancements
- **Job Management**
  - Executor supports special instance arguments

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
