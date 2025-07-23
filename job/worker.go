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
	"time"

	logger "github.com/cybergarage/go-logger/log"
)

// Worker is an interface that defines methods for processing jobs.
type Worker interface {
	// Start starts the worker to process jobs.
	Start() error
	// Stop cancels the worker from processing jobs.
	Stop() error
	// IsProcessing returns true if the worker is currently processing a job.
	IsProcessing() bool
	// StopWithWait stops the worker and waits for it to finish processing jobs.
	StopWithWait() error
}

type worker struct {
	queue      Queue
	done       chan struct{}
	processing bool
}

// IsProcessing returns true if the worker is currently processing a job.
func (w *worker) IsProcessing() bool {
	return w.processing
}

// WorkerOption is a function that configures a job worker.
type WorkerOption func(*worker)

// WithWorkerQueue sets the job queue for the worker.
func WithWorkerQueue(queue Queue) WorkerOption {
	return func(w *worker) {
		w.queue = queue
	}
}

// NewWorker creates a new instance of the job worker.
func NewWorker(opts ...WorkerOption) Worker {
	w := &worker{
		queue:      nil,
		done:       make(chan struct{}),
		processing: false,
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
		delay := ji.Policy().Backoff()
		if 0 < delay {
			time.Sleep(delay) // Wait for the retry delay
		}
		w.queue.Enqueue(ji) // Retry the job
		err := ji.UpdateState(JobScheduled)
		if err != nil {
			logError(ji, err)
		}
	}

	rescheduleInstance := func(ji Instance) {
		w.queue.Enqueue(ji) // Reschedule if recurring
		err := ji.UpdateState(JobScheduled)
		if err != nil {
			logError(ji, err)
		}
	}

	go func() {
		for {
			select {
			case <-w.done:
				w.processing = false
				return
			default:
				ji, err := w.queue.Dequeue()
				if err != nil {
					logger.Error(err)
					continue
				}
				err = ji.UpdateState(JobProcessing)
				if err != nil {
					logError(ji, err)
					continue
				}
				w.processing = true
				res, err := ji.Process()
				if err == nil {
					err = ji.UpdateState(JobCompleted, NewResultWith(res))
					if err != nil {
						logError(ji, err)
					}
					if ji.IsRecurring() {
						rescheduleInstance(ji)
					}
				} else {
					err = ji.UpdateState(JobTerminated, err)
					if err != nil {
						logError(ji, err)
					}
					if ji.IsRetriable() {
						retryInstance(ji)
					} else {
						if ji.IsRecurring() {
							rescheduleInstance(ji)
						}
					}
				}
				w.processing = false
			}
		}
	}()

	return nil
}

// StopWithWait stops the worker and waits for it to finish processing jobs.
func (w *worker) StopWithWait() error {
	for {
		if !w.IsProcessing() {
			break
		}
		time.Sleep(1 * time.Second) // Wait for worker to finish processing
	}
	return w.Stop()
}

// Stop cancels the worker from processing jobs.
func (w *worker) Stop() error {
	close(w.done)
	w.done = make(chan struct{})
	return nil
}
