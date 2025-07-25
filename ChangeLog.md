# ChangeLog

## 1.1.0 (2025-0X-XX)
- Job client
- Support 
  - Cancellation of job instances
  - Timeout handling for job instances
  - Executor supports instance arguments
- Prometheus
- CLI command
  - `go-jobctl` for job management
- Executor
  - Special arguments for job instances
- Clean before

## 1.0.0 (2025-0X-XX)
- gRPC API
- Distributed Store
  - etcd plugin added

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
