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
	"errors"
	"sync"
)

const (
	// DefaultWorkerNum is the default number of workers in the group.
	DefaultWorkerNum = 1
)

// WorkerGroup is an interface that defines methods for managing a group of workers.
type WorkerGroup interface {
	// Start starts all workers in the group.
	Start() error
	// Stop stops all workers in the group.
	Stop() error
	// StopWithWait stops all workers in the group and waits for them to finish processing.
	StopWithWait() error
	// ResizeWorkers scales the number of workers in the group.
	ResizeWorkers(num int) error
	// NumWorkers returns the number of workers in the group.
	NumWorkers() int
}

// WorkerGroupOption defines a function that configures a worker group.
type WorkerGroupOption func(*workerGroup)

// WithWorkerGroupNumber sets the number of workers in the group.
func WithNumWorkers(number int) WorkerGroupOption {
	return func(g *workerGroup) {
		g.workers = make([]Worker, number)
	}
}

// WithWorkerGroupQueue sets the job queue for the worker group.
func WithWorkerGroupQueue(queue Queue) WorkerGroupOption {
	return func(g *workerGroup) {
		g.queue = queue
	}
}

// NewWorkerGroup creates a new instance of the worker group with the provided options.
func NewWorkerGroup(opts ...WorkerGroupOption) WorkerGroup {
	return newWorkerGroup(opts...)
}

type workerGroup struct {
	sync.Mutex
	workers []Worker
	queue   Queue
}

func newWorkerGroup(opts ...WorkerGroupOption) *workerGroup {
	g := &workerGroup{
		Mutex:   sync.Mutex{},
		workers: make([]Worker, DefaultWorkerNum),
		queue:   nil,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Start starts all workers in the group.
func (g *workerGroup) Start() error {
	if g.queue == nil {
		return errors.New("worker group queue is not set")
	}
	for i := 0; i < len(g.workers); i++ {
		g.workers[i] = NewWorker(WithWorkerQueue(g.queue))
	}
	for _, w := range g.workers {
		if err := w.Start(); err != nil {
			return errors.Join(err, g.Stop())
		}
	}
	return nil
}

// Stop stops all workers in the group.
func (g *workerGroup) Stop() error {
	for i := 0; i < len(g.workers); i++ {
		if err := g.workers[i].Stop(); err != nil {
			return err
		}
	}
	return nil
}

// StopWithWait stops all workers in the group and waits for them to finish processing.
func (g *workerGroup) StopWithWait() error {
	for _, w := range g.workers {
		if err := w.StopWithWait(); err != nil {
			return err
		}
	}
	return nil
}

// ResizeWorkers scales the number of workers for the job manager.
func (g *workerGroup) ResizeWorkers(num int) error {
	if num <= 0 {
		return errors.New("number of workers must be positive")
	}
	if len(g.workers) == num {
		return nil
	}

	if !g.TryLock() {
		return errors.New("manager is scaling workers")
	}
	defer g.Unlock()

	if len(g.workers) < num {
		for i := len(g.workers); i < num; i++ {
			worker := NewWorker(WithWorkerQueue(g.queue))
			if err := worker.Start(); err != nil {
				return err
			}
			g.workers = append(g.workers, worker)
		}
	} else {
		for i := num; i < len(g.workers); i++ {
			if err := g.workers[i].StopWithWait(); err != nil {
				return err
			}
		}
		g.workers = g.workers[:num]
	}
	return nil
}

// NumWorkers returns the number of workers in the group.
func (g *workerGroup) NumWorkers() int {
	g.Lock()
	defer g.Unlock()
	return len(g.workers)
}
