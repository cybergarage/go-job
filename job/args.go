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

type Arguments interface {
	// Arguments returns the underlying arguments.
	Arguments() []any
}

type arguments struct {
	args []any
}

// ArgumentsOption defines a function that configures the arguments for a job.
type ArgumentsOption func(*arguments)

// WithArguments sets the arguments for a job.
func WithArguments(args ...any) ArgumentsOption {
	return func(a *arguments) {
		a.args = args
	}
}

func newArguments(opts ...ArgumentsOption) *arguments {
	a := &arguments{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Arguments returns the underlying arguments.
func (a *arguments) Arguments() []any {
	return a.args
}
