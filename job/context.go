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

// Context defines the interface for a job context that holds arguments
type Context interface {
	Arguments() []any
}

// ContextOption defines a function type that can be used to configure a job context.
type ContextOption func(*ctx)

// WithArguments sets the arguments for the job context.
func WithArguments(args ...any) ContextOption {
	return func(ctx *ctx) {
		ctx.args = args
	}
}

type ctx struct {
	args []any
}

// NewJobContext creates a new job context with the provided options.
func newJobContext(opts ...ContextOption) *ctx {
	ctx := &ctx{
		args: make([]any, 0),
	}
	for _, opt := range opts {
		opt(ctx)
	}
	return ctx
}

// Arguments returns the arguments for the job context.
func (c *ctx) Arguments() []any {
	return c.args
}
