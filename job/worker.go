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
	"context"
	"errors"
	"time"

	logger "github.com/cybergarage/go-logger/log"
)

// Worker is an interface that defines methods for processing jobs.
type Worker interface {
	// Start starts the worker to process jobs.
	Start() error
	// Cancel cancels the currently processing job. Returns an error if no job is being processed.
	Cancel() error
	// Wait waits for the worker to finish processing jobs.
	Wait() error
	// Stop cancels the worker from processing jobs.
	Stop() error
	// IsProcessing returns true if the worker is currently processing a job.
	IsProcessing() bool
	// ProcessingInstance returns the job instance being processed, if any.
	ProcessingInstance() (Instance, bool)
}

type worker struct {
	manager        Manager
	done           chan struct{}
	processingInst Instance
	jobCtx         context.Context
	jobCancel      context.CancelFunc
}

// IsProcessing returns true if the worker is currently processing a job.
func (w *worker) IsProcessing() bool {
	_, processing := w.ProcessingInstance()
	return processing
}

// ProcessingInstance returns the job instance being processed, if any.
func (w *worker) ProcessingInstance() (Instance, bool) {
	if w.processingInst == nil {
		return nil, false
	}
	return w.processingInst, true
}

// workerOption is a function that configures a job worker.
type workerOption func(*worker)

func withWorkerManager(mgr Manager) workerOption {
	return func(w *worker) {
		w.manager = mgr
	}
}

// newWorker creates a new instance of the job worker.
func newWorker(opts ...workerOption) Worker {
	w := &worker{
		manager:        nil,
		done:           make(chan struct{}),
		processingInst: nil,
		jobCtx:         nil,
		jobCancel:      nil,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Start starts the worker to process jobs.
func (w *worker) Start() error {
	if err := w.Stop(); err != nil {
		return err
	}
	return w.Run()
}

// Run starts the worker to process jobs in a loop.
func (w *worker) Run() error {
	logError := func(ji Instance, err error) {
		logger.Error(err)
		ji.Error(err)
	}
	retryInstance := func(ji Instance) {
		backoffStrategy := ji.Policy().BackoffStrategy()
		if backoffStrategy != nil {
			backoff := backoffStrategy(ji)
			if 0 < backoff {
				time.Sleep(backoff) // Wait for the retry delay
			}
		}
		w.manager.EnqueueInstance(ji) // Retry the job
		err := ji.UpdateState(JobScheduled)
		if err != nil {
			logError(ji, err)
		}
	}

	rescheduleInstance := func(ji Instance) {
		w.manager.EnqueueInstance(ji) // Reschedule if recurring
		err := ji.UpdateState(JobScheduled)
		if err != nil {
			logError(ji, err)
		}
	}

	go func() {
		for {
			select {
			case <-w.done:
				w.processingInst = nil
				return
			default:
				ji, err := w.manager.DequeueNextInstance()
				if err != nil {
					logger.Error(err)
					continue
				}
				err = ji.UpdateState(JobProcessing)
				if err != nil {
					logError(ji, err)
					continue
				}

				w.processingInst = ji

				w.jobCtx, w.jobCancel = context.WithCancel(context.Background())
				timeout := ji.Policy().Timeout()
				if 0 < timeout {
					w.jobCtx, w.jobCancel = context.WithTimeout(w.jobCtx, timeout)
				}

				res, err := ji.Process(w.jobCtx, w.manager, w, ji)

				w.jobCancel()
				w.jobCtx = nil
				w.jobCancel = nil

				if err == nil {
					err = ji.UpdateState(JobCompleted, newResultWith(res))
					if err != nil {
						logError(ji, err)
					}
					ji.CompleteProcessor()(ji, res)
					if ji.IsRecurring() {
						rescheduleInstance(ji)
					}
				} else {
					jobState := JobTerminated
					if errors.Is(err, context.Canceled) {
						jobState = JobCanceled
					} else if errors.Is(err, context.DeadlineExceeded) {
						jobState = JobTimedOut
					}
					err = ji.UpdateState(jobState, err)
					if err != nil {
						logError(ji, err)
					}
					if ji.IsRetriable() {
						retryInstance(ji)
					} else if ji.IsRecurring() {
						rescheduleInstance(ji)
					}
				}
				w.processingInst = nil
			}
		}
	}()

	return nil
}

// Wait waits for the worker to finish processing jobs.
func (w *worker) Wait() error {
	for w.IsProcessing() {
		time.Sleep(1 * time.Second) // Wait for worker to finish processing
	}
	return nil
}

// Cancel cancels the currently processing job. Returns an error if no job is being processed.
func (w *worker) Cancel() error {
	if !w.IsProcessing() {
		return ErrNotProcessing
	}
	if w.jobCancel != nil {
		w.jobCancel()
	}
	w.jobCtx = nil
	w.jobCancel = nil
	w.processingInst = nil
	return nil
}

// Stop cancels the worker from processing jobs.
func (w *worker) Stop() error {
	err := w.Cancel()
	if errors.Is(err, ErrNotProcessing) {
		// If not processing, just close the done channel
		err = nil
	}
	close(w.done)
	w.done = make(chan struct{})
	return err
}
