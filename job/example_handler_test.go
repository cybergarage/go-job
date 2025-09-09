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

package job_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/cybergarage/go-job/job"
)

func ExampleWithDescription() {
	// Create a job with a specific description
	job, _ := job.NewJob(
		job.WithKind("description"),
		job.WithExecutor(func() {}),
		job.WithDescription("This job sums two numbers"),
	)
	fmt.Printf("%s\n", job.Description())
	// Output: This job sums two numbers
}

func ExampleWithKind() {
	// Create a job with a specific kind
	job, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func() {}),
	)
	fmt.Printf("%s\n", job.Kind())
	// Output: sum
}

func ExampleWithExecutor_hello() {
	// Create a job with a specific executor
	job, _ := job.NewJob(
		job.WithKind("hello"),
		job.WithExecutor(func() {
			fmt.Println("Hello, world!")
		}),
	)
	fmt.Printf("%T\n", job.Handler().Executor())
	// Output: func()
}

func ExampleWithExecutor_sum() {
	// Create a job with a specific executor
	job, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func(a, b int) int {
			return a + b
		}),
	)
	fmt.Printf("%T\n", job.Handler().Executor())
	// Output: func(int, int) int
}

func ExampleWithExecutor_concat() {
	// Create a job with a specific executor
	job, _ := job.NewJob(
		job.WithKind("concat"),
		job.WithExecutor(func(a string, b string) string {
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
	job, _ := job.NewJob(
		job.WithKind("struct"),
		job.WithExecutor(func(param *ConcatString) *ConcatString {
			param.S = param.A + ", " + param.B
			return param
		}),
	)
	fmt.Printf("%T\n", job.Handler().Executor())

	// Output:
	// func(*job_test.ConcatString) *job_test.ConcatString
}

func ExampleWithCompleteProcessor() {
	// Create a job with a specific complete processor
	job, _ := job.NewJob(
		job.WithKind("complete"),
		job.WithExecutor(func() {}),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			ji.Infof("Job completed with result: %v", res)
		}),
	)
	fmt.Printf("%T\n", job.Handler().CompleteProcessor())
	// Output: job.CompleteProcessor
}

func ExampleWithTerminateProcessor() {
	// Create a job with a specific terminate processor
	job, _ := job.NewJob(
		job.WithKind("terminate"),
		job.WithExecutor(func() {}),
		job.WithTerminateProcessor(func(ji job.Instance, err error) error {
			if errors.Is(err, context.DeadlineExceeded) {
				// Do not retry if the job was terminated due to a deadline being exceeded
				ji.Infof("Job (%s) terminated due to deadline exceeded: %v", ji.Kind(), err)
				return nil
			}
			// Retry for all other errors
			return err
		}),
	)
	fmt.Printf("%T\n", job.Handler().TerminateProcessor())
	// Output: func(job.Instance, error) error
}

func ExampleWithStateChangeProcessor() {
	// Create a job with a specific state change processor
	job, _ := job.NewJob(
		job.WithKind("state"),
		job.WithExecutor(func() {}),
		job.WithStateChangeProcessor(func(ji job.Instance, state job.JobState) {
			ji.Infof("State changed to: %v", state)
		}),
	)
	fmt.Printf("%T\n", job.Handler().StateChangeProcessor())
	// Output: func(job.Instance, job.JobState)
}
