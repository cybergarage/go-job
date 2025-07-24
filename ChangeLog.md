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

## v0.9.2 (2025-07-XX)
### 🚀 New Features
- **Instance Querying & Management**
  - Added `Manager::ListInstances()` and `LookupInstances()` for comprehensive instance management
  - Introduced `Query` system to lookup scheduled and completed instances with filtering capabilities
  - Added `NewInstancesFromQueue()` for creating instances from queue data
  - Enhanced `StateHistory::ListInstanceHistory()` to retrieve complete state history
- **Enhanced Data Structures**
  - Added `InstanceMap::Arguments()`, `ResultSet()`, and `Error()` methods for better data access
  - Introduced `NewArgumentsFrom()` and `Arguments::JSONString()` for improved argument handling
  - Added `NewInstanceMapWith()` for flexible map-based instance creation
  - Implemented map key constants for consistent data access
- **Job State Management**
  - Added `JobStateAll` constant for comprehensive state handling
  - Renamed `JobStateUnknown` to `JobStateUnset` for better clarity
  - Enhanced state tracking with proper timestamp management
### 🔧 Improvements
- **Instance Lifecycle**
  - Updated `NewInstance()` to handle job options more effectively
  - Enhanced `NewInstancesFromHistory()` to:
    - Parse and store argument options
    - Set schedule, processed, and created timestamps correctly
    - Handle instance result sets and errors properly
    - Store argument options with improved data consistency
- **Manager Enhancements**
  - Updated `Manager::ScheduleJob()` to properly update created status
  - Improved scheduling with better timestamp handling
  - Enhanced job registration and lifecycle management
- **Worker & Processing**
  - Updated `Worker::Run()` to:
    - Store errors when jobs are terminated
    - Update job status during rescheduling
    - Better handle job completion and termination states
- **Queue Management**
  - Updated `Queue::Clear()` functionality
  - Added `Queue::List()` to return all queued instances
  - Removed `Queue::HasJobs()` in favor of using `Empty()`
### 📝 API & Protocol Updates
- **Job Interface**
  - Added `Job::Description()` and `WithDescription()` for job metadata
  - Added `Job::RegisteredAt()` to track registration timestamps
  - Updated `Registry::ListJobs()` to return proper error handling
- **Instance Interface**
  - Added `Instance::ProcessedAt()` for execution tracking
  - Enhanced `Instance::CreatedAt()`, `CompletedAt()`, and `TerminatedAt()` timestamp methods
  - Added `Instance::ResultSet()` to return processed results for completed/terminated jobs
  - Added `WithUUID()` option to pass UUID to instances
- **Protocol Buffers**
  - Multiple updates to `service.proto` for better API structure
  - Enhanced message definitions for improved client-server communication
### 🛠️ Internal Improvements
- **Data Handling**
  - Updated `encoding.UnmarshalToMap()` and `UnmarshalJSONToMap()` with proper error returns
  - Enhanced `Policy::Map()` for better policy serialization
  - Improved `InstanceMap::Map()` functionality
- **Testing & Quality**
  - Updated `TestScheduleJobs()` with more comprehensive test coverage
  - Enhanced `TestArgumentsFrom()` and `TestJobState()` for better validation
  - Updated `TestSchedules()` for improved scheduling tests
- **Build & Maintenance**
  - Updated go.mod dependencies
  - Enhanced Makefile for better build process
  - Updated documentation images
### 🔄 Breaking Changes
- **State Naming**: `JobStateUnknown` has been renamed to `JobStateUnset`
- **Result Interface**: Renamed `Result` to `ResultSet` for better clarity
### 📋 Bug Fixes
- Fixed timestamp handling in various instance operations
- Improved error handling in encoding/decoding operations
- Enhanced data consistency in instance history management
- Better handling of job termination and completion states

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
