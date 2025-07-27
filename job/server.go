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
	"os"

	logger "github.com/cybergarage/go-logger/log"
)

// Server is an interface that defines methods for managing the job server.
type Server interface {
	// Manager returns the job manager associated with the server.
	Manager() Manager
	// Start starts the job server.
	Start() error
	// Stop stops the job server.
	Stop() error
}

type server struct {
	manager Manager
}

// NewServer returns a new job server instance.
func NewServer(opts ...any) (Server, error) {
	mgr, err := NewManager(opts...)
	if err != nil {
		return nil, err
	}
	return &server{
		manager: mgr,
	}, nil
}

// Manager returns the job manager associated with the server.
func (server *server) Manager() Manager {
	return server.manager
}

// Start starts the job server.
func (server *server) Start() error {
	starters := []func() error{
		server.manager.Start,
	}
	var errs error
	for _, starter := range starters {
		if err := starter(); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	logger.Infof("%s (PID:%d) started", ProductName, os.Getpid())

	return errs
}

// Stop stops the job server.
func (server *server) Stop() error {
	stoppers := []func() error{
		server.manager.Stop,
	}
	var errs error
	for _, stopper := range stoppers {
		if err := stopper(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	logger.Infof("%s (PID:%d) terminated", ProductName, os.Getpid())
	return errs
}
