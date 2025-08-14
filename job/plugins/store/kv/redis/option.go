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

package redis

import (
	"net"

	redis "github.com/redis/go-redis/v9"
)

const (
	// DefaultPort is the default port for the Redis store.
	DefaultPort = "6379"

	// DefaultHost is the default host for the Redis store.
	DefaultHost = "127.0.0.1"
)

// StoreOption is an alias for redis.Options, used for configuring the Redis store.
type StoreOption = redis.Options

// NewStoreOption creates a new StoreOption with the specified options.
func NewStoreOption(opts ...any) StoreOption {
	defaultAddr := net.JoinHostPort(DefaultHost, DefaultPort)
	// nolint:exhaustruct
	sopt := redis.Options{
		Addr: defaultAddr,
	}
	for _, opt := range opts {
		switch v := opt.(type) {
		case net.Addr:
			sopt.Addr = v.String()
		}
	}
	return sopt
}
