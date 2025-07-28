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

import (
	"fmt"
)

// Registry is an interface that defines methods for managing job instances.
type Registry interface {
	// RegisterJob registers a job in the registry.
	RegisterJob(job Job) error
	// UnregisterJob removes a job from the registry by its kind.
	UnregisterJob(kind Kind) error
	// ListJobs returns a slice of all registered jobs.
	ListJobs() ([]Job, error)
	// LookupJob looks up a job by its kind in the registry.
	LookupJob(kind Kind) (Job, bool)
	// Clear clears all registered jobs.
	Clear() error
}

// registry is responsible for managing job instances.
type registry struct {
	jobs map[string]Job
}

// NewRegistry creates a new instance of Registry.
func NewRegistry() Registry {
	return &registry{
		jobs: make(map[string]Job),
	}
}

// RegisterJob registers a job in the registry.
func (reg *registry) RegisterJob(job Job) error {
	if _, exists := reg.jobs[job.Kind()]; exists {
		return fmt.Errorf("job with kind %q is already registered", job.Kind())
	}
	reg.jobs[job.Kind()] = job
	return nil
}

// UnregisterJob removes a job from the registry by its kind.
func (reg *registry) UnregisterJob(kind Kind) error {
	if _, exists := reg.jobs[kind]; !exists {
		return fmt.Errorf("job with kind %q is not registered", kind)
	}
	delete(reg.jobs, kind)
	return nil
}

// ListJobs returns a slice of all registered jobs.
func (reg *registry) ListJobs() ([]Job, error) {
	jobs := make([]Job, 0, len(reg.jobs))
	for _, job := range reg.jobs {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// LookupJob looks up a job by its kind in the registry.
func (reg *registry) LookupJob(kind Kind) (Job, bool) {
	job, exists := reg.jobs[kind]
	return job, exists
}

// Clear clears all registered jobs.
func (reg *registry) Clear() error {
	reg.jobs = make(map[string]Job)
	return nil
}
