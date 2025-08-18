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

import "fmt"

func ExampleWithDescription() {
	// Create a job with a specific description
	job, _ := NewJob(
		WithDescription("This job sums two numbers"),
	)
	fmt.Printf("%s\n", job.Description())
	// Output: This job sums two numbers
}

func ExampleWithKind() {
	// Create a job with a specific kind
	job, _ := NewJob(
		WithKind("sum"),
	)
	fmt.Printf("%s\n", job.Kind())
	// Output: sum
}

func ExampleWithExecutor_hello() {
	// Create a job with a specific executor
	job, _ := NewJob(
		WithExecutor(func() {
			fmt.Println("Hello, world!")
		}),
	)
	fmt.Printf("%T\n", job.Handler().Executor())
	// Output: func()
}

func ExampleWithExecutor_sum() {
	// Create a job with a specific executor
	job, _ := NewJob(
		WithExecutor(func(a, b int) int {
			return a + b
		}),
	)
	fmt.Printf("%T\n", job.Handler().Executor())
	// Output: func(int, int) int
}

func ExampleWithExecutor_concat() {
	// Create a job with a specific executor
	job, _ := NewJob(
		WithExecutor(func(a string, b string) string {
			return a + ", " + b
		}),
	)
	fmt.Printf("%T\n", job.Handler().Executor())
	// Output: func(string, string) string
}

func ExampleWithExecutor_struct() {
	type ConcatString struct {
		A string
		B string
		S string
	}
	job, _ := NewJob(
		WithExecutor(func(param *ConcatString) *ConcatString {
			param.S = param.A + ", " + param.B
			return param
		}),
	)
	fmt.Printf("%T\n", job.Handler().Executor())

	// Output:
	// func(*job.ConcatString) *job.ConcatString
}
