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
	"encoding/json"
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

// NewArgumentsWith creates a new Arguments instance with the provided arguments.
func NewArgumentsWith(args ...any) Arguments {
	return newArguments(WithArguments(args...))
}

// NewArgumentsWithStrings creates a new Arguments instance with the provided string arguments.
// It converts the strings to any type for compatibility with the Arguments interface.
func NewArgumentsWithStrings(args ...string) Arguments {
	anyArgs := make([]any, len(args))
	for i, arg := range args {
		anyArgs[i] = arg
	}
	return newArguments(WithArguments(anyArgs...))
}

// NewArgumentsWithString creates a new Arguments instance from a JSON string representation of arguments.
func NewArgumentsWithString(s string) (Arguments, error) {
	var arr []any
	err := json.Unmarshal([]byte(s), &arr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}
	return NewArgumentsWith(arr...), nil
}

// NewArgumentsFrom creates a new Arguments instance from the provided arguments.
// It supports various types such as nil, []any, and string.
func NewArgumentsFrom(args any) (Arguments, error) {
	switch v := args.(type) {
	case nil:
		return NewArguments(), nil
	case Arguments:
		return v, nil
	case []any:
		return NewArgumentsWith(v...), nil
	case []string:
		return NewArgumentsWithStrings(v...), nil
	case string:
		return NewArgumentsWithString(v)
	}
	return nil, fmt.Errorf("unsupported type for arguments: %T", args)
}

// NewArguments creates a new Arguments instance with the provided options.
// It allows for flexible configuration of the arguments.
func NewArguments(opts ...ArgumentsOption) Arguments {
	return newArguments(opts...)
}

func newArguments(opts ...ArgumentsOption) *arguments {
	args := &arguments{
		Args: make([]any, 0),
	}
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
		argumentsKey: args.String(),
	}
}

// String returns a string representation of the arguments.
func (args *arguments) String() string {
	return fmt.Sprintf("%v", args.Args)
}
