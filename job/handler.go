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

// StateChangeProcessor is called when a job's state changes.
type StateChangeProcessor = func(job Instance, state JobState)

// TerminateProcessor is called when a job reaches the terminated state.
// Users can:
//   - Return nil to resolve the error (mark job as successful)
//   - Return a modified error to transform the error
//   - Return the original error to keep it unchanged
type TerminateProcessor = func(job Instance, err error) error

// CompleteProcessor is called when a job reaches the completed state (successful completion).
// It allows users to handle the results of the job execution.
// Users can process the results and perform any necessary actions.
type CompleteProcessor func(job Instance, responses []any)

// HandlerOption is a function type that applies options to a job handler.
type HandlerOption func(*handler)

// WithExecutor sets the executor function for the job handler.
func WithExecutor(executor Executor) HandlerOption {
	return func(h *handler) {
		h.executor = executor
	}
}

// WithStateChangeProcessor sets a handler function that is invoked each time the state of a job instance changes while being processed by the local worker.
// NOTE: In a distributed environment with multiple worker groups, the worker that schedules a job instance may not receive all status updates for that instance.
func WithStateChangeProcessor(fn StateChangeProcessor) HandlerOption {
	return func(h *handler) {
		h.stateChgProcessor = fn
	}
}

// WithTerminateProcessor sets a handler function that is called if a job instance ends with an error during execution by the local worker.
// NOTE: In a distributed environment with multiple worker groups, the worker that schedules a job instance may be different from the worker that actually executes it.
func WithTerminateProcessor(fn TerminateProcessor) HandlerOption {
	return func(h *handler) {
		h.terminateProcessor = fn
	}
}

// WithCompleteProcessor sets a handler function that is called when a job instance completes successfully during execution by the local worker.
// NOTE: In a distributed environment with multiple worker groups, the worker that schedules a job instance may be different from the worker that actually executes it.
func WithCompleteProcessor(fn CompleteProcessor) HandlerOption {
	return func(h *handler) {
		h.completeProcessor = fn
	}
}

// Handler is an interface that defines methods for executing jobs and handling errors.
type Handler interface {
	// Executor returns the executor function set for the job handler.
	Executor() Executor
	// StateChangeProcessor returns the state change handler function set for the job handler.
	StateChangeProcessor() StateChangeProcessor
	// CompleteProcessor returns the completion handler function set for the job handler.
	CompleteProcessor() CompleteProcessor
	// TerminateProcessor returns the error processor function set for the job handler.
	TerminateProcessor() TerminateProcessor
	// Execute runs the job with the provided parameters.
	Execute(params ...any) ([]any, error)
	// HandleTerminated processes errors that occur during job execution.
	HandleTerminated(job Instance, err error) error
	// HandleCompleted processes the responses from a job execution.
	HandleCompleted(job Instance, responses []any)
}

type handler struct {
	executor           Executor
	stateChgProcessor  StateChangeProcessor
	terminateProcessor TerminateProcessor
	completeProcessor  CompleteProcessor
}

// NewHandler creates a new instance of a job handler with the provided options.
func NewHandler(opts ...HandlerOption) Handler {
	return newHandler(opts...)
}

func newHandler(opts ...HandlerOption) *handler {
	h := &handler{
		executor:           nil,
		stateChgProcessor:  nil,
		terminateProcessor: nil,
		completeProcessor:  nil,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Executor returns the executor function set for the job handler.
func (h *handler) Executor() Executor {
	return h.executor
}

// StateChangeProcessor returns the state change handler function set for the job handler.
func (h *handler) StateChangeProcessor() StateChangeProcessor {
	return h.stateChgProcessor
}

// TerminateProcessor returns the error processor function set for the job handler.
func (h *handler) TerminateProcessor() TerminateProcessor {
	return h.terminateProcessor
}

// CompleteProcessor returns the response handler function set for the job handler.
func (h *handler) CompleteProcessor() CompleteProcessor {
	return h.completeProcessor
}

// Execute runs the job using the executor function, if set.
func (h *handler) Execute(params ...any) ([]any, error) {
	if h.executor == nil {
		return nil, fmt.Errorf("no executor set for job handler")
	}
	res, err := Execute(h.executor, params...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// HandleTerminated processes errors that occur during job execution using the error handler, if set.
func (h *handler) HandleTerminated(job Instance, err error) error {
	if h.terminateProcessor == nil {
		return err
	}
	return h.terminateProcessor(job, err)
}

// HandleCompleted processes the responses from a job execution using the response handler, if set.
func (h *handler) HandleCompleted(job Instance, responses []any) {
	if h.completeProcessor == nil {
		return
	}
	h.completeProcessor(job, responses)
}
