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

// JobErrorHandler is a function type that defines how to handle errors during job execution.
type JobErrorHandler = func(job Job, err error) (bool, error)

// JobHandlerOption is a function type that applies options to a job handler.
type JobHandlerOption func(*jobHandler)

// WithExecutor sets the executor function for the job handler.
func WithExecutor(executor JobExecutor) JobHandlerOption {
	return func(h *jobHandler) {
		h.executor = executor
	}
}

// WithErrorHandler sets the error handler function for the job handler.
func WithErrorHandler(errorHandler JobErrorHandler) JobHandlerOption {
	return func(h *jobHandler) {
		h.errorHandler = errorHandler
	}
}

// JobHandler is an interface that defines methods for executing jobs and handling errors.
type JobHandler interface {
	// Execute runs the job with the provided parameters.
	Execute(params ...any) error
	// HandleError processes errors that occur during job execution.
	HandleError(job Job, err error) (bool, error)
}

type jobHandler struct {
	executor     JobExecutor
	errorHandler JobErrorHandler
}

func newJobHandler(opts ...JobHandlerOption) *jobHandler {
	h := &jobHandler{
		executor:     nil, // Default executor is nil
		errorHandler: nil, // Default error handler is nil
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Execute runs the job using the executor function, if set.
func (h *jobHandler) Execute(params ...any) error {
	if h.executor == nil {
		return nil // No executor set, nothing to do
	}
	_, err := ExecuteJob(h.executor, params...)
	if err != nil {
		return err
	}
	return nil
}

func (h *jobHandler) HandleError(job Job, err error) (bool, error) {
	if h.errorHandler == nil {
		return false, err // No error handler set, return the error
	}
	return h.errorHandler(job, err)
}
