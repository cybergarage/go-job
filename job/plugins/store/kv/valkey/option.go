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

package valkey

import (
	"net"

	"github.com/valkey-io/valkey-go"
)

const (
	// DefaultPort is the default port for the Valkey store.
	DefaultPort = "6379"

	// DefaultHost is the default host for the Valkey store.
	DefaultHost = "127.0.0.1"
)

// StoreOption is an alias for valkey.ClientOption, used for configuring the Valkey store.
type StoreOption = valkey.ClientOption

// NewStoreOption creates a new StoreOption with the specified options.
// nolint:exhaustruct
func NewStoreOption(opts ...any) StoreOption {
	defaultAddr := net.JoinHostPort(DefaultHost, DefaultPort)
	opt := valkey.ClientOption{
		InitAddress: []string{defaultAddr},
	}
	for _, o := range opts {
		switch v := o.(type) {
		case net.Addr:
			opt.InitAddress = []string{v.String()}
		}
	}
	return opt
}
