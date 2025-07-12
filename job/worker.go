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

// Worker is an interface that defines methods for processing jobs.
type Worker interface {
	// Start starts the worker to process jobs.
	Start() error
	// Stop stops the worker from processing jobs.
	Stop() error
}

type worker struct {
	queue Queue
	done  chan struct{}
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
		queue: nil,
		done:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Start starts the worker to process jobs.
func (w *worker) Start() error {
	go func() {
		for {
			select {
			case <-w.done:
				return
			default:
				ji, err := w.queue.Dequeue()
				if err != nil {
					continue
				}
				err = ji.UpdateState(JobProcessing)
				if err != nil {
					continue
				}
				err = ji.Process()
				if err != nil {
					continue
				}
			}
		}
	}()
	return nil
}

// Stop stops the worker from processing jobs.
func (w *worker) Stop() error {
	close(w.done)
	return nil
}
