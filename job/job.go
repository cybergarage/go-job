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

// Job represents a job that can be scheduled to run at a specific time or interval.
type Job interface {
	// Kind returns the name of the job.
	Kind() string
	// Handler returns the job handler for the job.
	Handler() JobHandler
	// Payload returns the payload of the job.
	Payload() any
	// Process processes the job using the job handler.
	Process() error
	// String returns a string representation of the job.
	String() string
}
