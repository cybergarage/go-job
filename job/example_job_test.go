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
	"math"
	"strings"
)

func ExampleNewJob_simple() {
	job, err := NewJob(
		WithKind("no args and no return"),
		WithDescription("A simple job that prints a message"),
		WithExecutor(func() { fmt.Println("Hello, World!") }),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())
	// Output: Created job: no args and no return
}

func ExampleNewJob_concat() {
	job, err := NewJob(
		WithKind("concat (two args and one return)"),
		WithDescription("Concatenates two strings"),
		WithExecutor(func(a, b string) string { return a + ", " + b }),
		WithCompleteProcessor(func(ji Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())
	// Output: Created job: concat (two args and one return)
}

func ExampleNewJob_split() {
	job, err := NewJob(
		WithKind("split (one arg and two return)"),
		WithDescription("Splits a string into two parts"),
		WithExecutor(func(s string) (string, string) {
			parts := strings.Split(s, ",")
			return parts[0], parts[1]
		}),
		WithCompleteProcessor(func(ji Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0], res[1])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())
	// Output: Created job: split (one arg and two return)
}

func ExampleNewJob_abs() {
	job, err := NewJob(
		WithKind("abs (one arg and one return)"),
		WithDescription("Returns the absolute value of an integer"),
		WithExecutor(func(a int) int { return int(math.Abs(float64(a))) }),
		WithCompleteProcessor(func(ji Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())
	// Output: Created job: abs (one arg and one return)
}

func ExampleNewJob_sum() {
	job, err := NewJob(
		WithKind("sum (two args and one return)"),
		WithDescription("Returns the sum of two integers"),
		WithExecutor(func(a, b int) int { return a + b }),
		WithCompleteProcessor(func(ji Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())
	// Output: Created job: sum (two args and one return)
}

func ExampleNewJob_struct() {
	type SumOpt struct {
		a int
		b int
	}
	job, err := NewJob(
		WithKind("sum (struct arg and one return)"),
		WithDescription("Returns the sum of two integers from a struct"),
		WithExecutor(func(opt SumOpt) int { return opt.a + opt.b }),
		WithCompleteProcessor(func(ji Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())
	// Output: Created job: sum (struct arg and one return)
}

func ExampleNewJob_mutatingStruct() {
	type concatString struct {
		a string
		b string
		s string
	}
	job, err := NewJob(
		WithKind("concat (one struct input and one struct output)"),
		WithDescription("Concatenates two strings from a struct"),
		WithExecutor(func(param *concatString) *concatString {
			// Store the concatenated string result in the input struct
			param.s = param.a + " " + param.b
			return param
		}),
		WithCompleteProcessor(func(ji Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())
	// Output: Created job: concat (one struct input and one struct output)
}
