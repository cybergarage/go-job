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

// TerminateProcessor defines how users can process job execution errors.
// Users can:
//   - Return nil to resolve the error (mark job as successful)
//   - Return a modified error to transform the error
//   - Return the original error to keep it unchanged
type TerminateProcessor = func(job Instance, err error) error

// ResponseHandler is a function type that handles the responses from a job execution.
type ResponseHandler func(job Instance, responses []any)

// HandlerOption is a function type that applies options to a job handler.
type HandlerOption func(*handler)

// WithExecutor sets the executor function for the job handler.
func WithExecutor(executor Executor) HandlerOption {
	return func(h *handler) {
		h.executor = executor
	}
}

// WithTerminateProcessor sets the error handler function for the job handler.
func WithTerminateProcessor(fn TerminateProcessor) HandlerOption {
	return func(h *handler) {
		h.errorProcessor = fn
	}
}

// WithResponseHandler sets the response handler function for the job handler.
func WithResponseHandler(fn ResponseHandler) HandlerOption {
	return func(h *handler) {
		h.resHandler = fn
	}
}

// Handler is an interface that defines methods for executing jobs and handling errors.
type Handler interface {
	// Executor returns the executor function set for the job handler.
	Executor() Executor
	// TerminateProcessor returns the error processor function set for the job handler.
	TerminateProcessor() TerminateProcessor
	// ResponseHandler returns the response handler function set for the job handler.
	ResponseHandler() ResponseHandler
	// Execute runs the job with the provided parameters.
	Execute(params ...any) ([]any, error)
	// HandleError processes errors that occur during job execution.
	HandleError(job Instance, err error) error
	// HandleResponse processes the responses from a job execution.
	HandleResponse(job Instance, responses []any)
}

type handler struct {
	executor       Executor
	errorProcessor TerminateProcessor
	resHandler     ResponseHandler
}

// NewHandler creates a new instance of a job handler with the provided options.
func NewHandler(opts ...HandlerOption) Handler {
	return newHandler(opts...)
}

func newHandler(opts ...HandlerOption) *handler {
	h := &handler{
		executor:       nil,
		errorProcessor: nil,
		resHandler:     nil,
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

// TerminateProcessor returns the error processor function set for the job handler.
func (h *handler) TerminateProcessor() TerminateProcessor {
	return h.errorProcessor
}

// ResponseHandler returns the response handler function set for the job handler.
func (h *handler) ResponseHandler() ResponseHandler {
	return h.resHandler
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

// HandleError processes errors that occur during job execution using the error handler, if set.
func (h *handler) HandleError(job Instance, err error) error {
	if h.errorProcessor == nil {
		return err
	}
	return h.errorProcessor(job, err)
}

// HandleResponse processes the responses from a job execution using the response handler, if set.
func (h *handler) HandleResponse(job Instance, responses []any) {
	if h.resHandler == nil {
		return
	}
	h.resHandler(job, responses)
}
