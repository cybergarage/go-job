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

type Arguments interface {
	// Arguments returns the underlying arguments.
	Arguments() []any
	// Map returns the arguments as a map.
	Map() map[string]any
	// String returns a string representation of the arguments.
	String() string
}

type arguments struct {
	Args []any `json:"args"`
}

// ArgumentsOption defines a function that configures the arguments for a job.
type ArgumentsOption func(*arguments)

// WithArguments sets the arguments for a job.
func WithArguments(args ...any) ArgumentsOption {
	return func(a *arguments) {
		a.Args = args
	}
}

func newArguments(opts ...ArgumentsOption) *arguments {
	args := &arguments{}
	for _, opt := range opts {
		opt(args)
	}
	return args
}

// Arguments returns the underlying arguments.
func (args *arguments) Arguments() []any {
	return args.Args
}

// Map returns the arguments as a map.
func (args *arguments) Map() map[string]any {
	return map[string]any{
		argsKey: args.String(),
	}
}

// String returns a string representation of the arguments.
func (args *arguments) String() string {
	return fmt.Sprintf("%v", args.Args)
}
