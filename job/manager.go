// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package job

// JobManager is an interface that defines methods for managing jobs.
type JobManager interface {
	// JobRegistry provides access to the job registry.
	JobRegistry
	// JobScheduler provides access to the job scheduler.
	JobScheduler
	// Start starts the job manager.
	Start() error
	// Stop stops the job manager.
	Stop() error
}
