// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"context"
	"net"
	"time"

	v3 "go.etcd.io/etcd/client/v3"
)

const (
	// DefaultPort is the default port for the etcd store.
	DefaultPort = "2379"

	// DefaultHost is the default host for the etcd store.
	DefaultHost = "127.0.0.1"

	// DefaultDialTimeout is the default dial timeout for the etcd store.
	DefaultDialTimeout = 5 * time.Second
)

// StoreOption is an alias for v3.Config, used for configuring the etcd store.
type StoreOption = v3.Config

// NewStoreOption creates a new StoreOption with the default values.
func NewStoreOption(opts ...any) StoreOption {
	defaultAddr := net.JoinHostPort(DefaultHost, DefaultPort)
	// nolint:exhaustruct
	sopt := v3.Config{
		Context:     context.Background(),
		Endpoints:   []string{defaultAddr},
		DialTimeout: DefaultDialTimeout,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case net.Addr:
			sopt.Endpoints = []string{v.String()}
		}
	}
	return sopt
}
