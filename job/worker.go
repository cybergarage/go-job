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
					logger.Error(err)
					continue
				}
				w.processing = true
				_, err = ji.Process()
				if err == nil {
					err = ji.UpdateState(JobCompleted)
					if err != nil {
						logger.Error(err)
					}
				} else {
					err = ji.UpdateState(JobTerminated)
					if err != nil {
						logger.Error(err)
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
