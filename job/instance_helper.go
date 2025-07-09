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

// InstanceHelper defines methods that can be used to perform actions before and after processing a job instance.
type InstanceHelper interface {
	// Before checks if the job instance should be processed before the given instance.
	Before(Instance) bool
	// After checks if the job instance should be processed after the given instance.
	After(Instance) bool
}

// Before checks if the job instance should be processed before the given instance.
func (ji *jobInstance) Before(other Instance) bool {
	if other == nil {
		return false
	}
	if ji.Policy().Priority() != other.Policy().Priority() {
		return ji.Policy().Priority() > other.Policy().Priority()
	}
	return ji.ScheduledAt().Before(other.ScheduledAt())
}

// After checks if the job instance should be processed after the given instance.
func (ji *jobInstance) After(other Instance) bool {
	return !ji.Before(other)
}
