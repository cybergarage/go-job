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

// Arguments represents the arguments passed to a job.
type Arguments struct {
	args []any
}

// ArgumentsOption defines a function that configures the arguments for a job.
type ArgumentsOption func(*Arguments)

// WithArguments sets the arguments for a job.
func WithArguments(args ...any) ArgumentsOption {
	return func(a *Arguments) {
		a.args = args
	}
}

// NewArguments creates a new Arguments instance.
func NewArguments(opts ...ArgumentsOption) *Arguments {
	a := &Arguments{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Args returns the underlying arguments.
func (a *Arguments) Args() []any {
	return a.args
}
