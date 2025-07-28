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

// Client represents a gRPC client.
type Client interface {
	// Name returns the name of the client.
	Name() string
	// SetHost sets a host name.
	SetHost(host string)
	// SetPort sets a port number.
	SetPort(port int)
	// Open opens a connection.
	Open() error
	// Close closes the connection.
	Close() error
	// GetVersion retrieves the version of the service.
	GetVersion() (string, error)
	// ScheduleJob schedules a job with the given kind and arguments.
	ScheduleJob(kind string, args ...any) (Instance, error)
	// ListRegisteredJobs lists all registered jobs.
	ListRegisteredJobs() ([]Job, error)
	// LookupInstances looks up job instances based on the provided query.
	LookupInstances(query Query) ([]Instance, error)
}

// NewClient returns a new default gRPC client.
func NewClient() Client {
	return NewGrpcClient()
}
