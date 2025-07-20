# ChangeLog

## 1.1.0 (2025-0X-XX)
- Job client
- Support 
  - Cancellation of job instances
  - Timeout handling for job instances
  - Executor supports instance arguments

## 1.0.1 (2025-0X-XX)
- Executor
  - Special arguments for job instances
- Clean before

## 1.0.0 (2025-0X-XX)
- gRPC API
- Distributed Store
  - etcd plugin added
- CLI command
  - `go-jobctl` for job management

## v0.9.1 (2025-07-19)
- Logging more system errors 
- Histrory and Logging cleanup
- Lookup logs and history by job ID

## v0.9.1 (2025-07-190)
### ✨ Features
- Added `Clear()` method to `Repository`
- Added `ClearInstanceLogs()` to `LogStore`
- Added `ClearInstanceHistory()` to `StateStore`
- Added `ClearInstances()` to `QueueStore` interface
- Extended `Instance` interface with `Logs()` method
### 🛠 Enhancements
- Extended `Log` interface with `Kind()` method
- Enabled backoff strategy in `Worker::Run()`
- Improved backoff policy interface
- Enhanced `Worker::Run()` to log job errors

## v0.9.0 (2025-07-17)
- Initial release
### Major Features
- Flexible job creation and registration with custom executors and arguments
- Priority and policy management for job execution order
- Job execution with response and error handlers
- Job state and history logging for monitoring and debugging
