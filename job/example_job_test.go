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
	"fmt"
	"math"
	"strings"

	"github.com/cybergarage/go-job/job"
)

func ExampleNewJob_simple() {
	job, err := job.NewJob(
		job.WithKind("no args and no return"),
		job.WithDescription("A simple job that prints a message"),
		job.WithExecutor(func() { fmt.Println("Hello, World!") }),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())

	// Output:
	// Created job: no args and no return
}

func ExampleNewJob_concat() {
	job, err := job.NewJob(
		job.WithKind("concat (two args and one return)"),
		job.WithDescription("Concatenates two strings"),
		job.WithExecutor(func(a, b string) string { return a + ", " + b }),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())

	// Output:
	// Created job: concat (two args and one return)
}

func ExampleNewJob_split() {
	job, err := job.NewJob(
		job.WithKind("split (one arg and two return)"),
		job.WithDescription("Splits a string into two parts"),
		job.WithExecutor(func(s string) (string, string) {
			parts := strings.Split(s, ",")
			return parts[0], parts[1]
		}),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0], res[1])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())

	// Output:
	// Created job: split (one arg and two return)
}

func ExampleNewJob_abs() {
	job, err := job.NewJob(
		job.WithKind("abs (one arg and one return)"),
		job.WithDescription("Returns the absolute value of an integer"),
		job.WithExecutor(func(a int) int { return int(math.Abs(float64(a))) }),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())

	// Output:
	// Created job: abs (one arg and one return)
}

func ExampleNewJob_sum() {
	job, err := job.NewJob(
		job.WithKind("sum (two args and one return)"),
		job.WithDescription("Returns the sum of two integers"),
		job.WithExecutor(func(a, b int) int { return a + b }),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())

	// Output:
	// Created job: sum (two args and one return)
}

func ExampleNewJob_struct() {
	type SumOpt struct {
		A int
		B int
	}
	job, err := job.NewJob(
		job.WithKind("sum (struct arg and one return)"),
		job.WithDescription("Returns the sum of two integers from a struct"),
		job.WithExecutor(func(opt SumOpt) int { return opt.A + opt.B }),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())

	// Output:
	// Created job: sum (struct arg and one return)
}

func ExampleNewJob_mutatingStruct() {
	type ConcatString struct {
		A string
		B string
		S string
	}
	job, err := job.NewJob(
		job.WithKind("concat (one struct input and one struct output)"),
		job.WithDescription("Concatenates two strings from a struct"),
		job.WithExecutor(func(param *ConcatString) *ConcatString {
			// Store the concatenated string result in the input struct
			param.S = param.A + " " + param.B
			return param
		}),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			// In this case, log the result to the go-job manager
			ji.Infof("%v", res[0])
		}),
	)
	if err != nil {
		fmt.Printf("Error creating job: %v\n", err)
		return
	}
	fmt.Printf("Created job: %s\n", job.Kind())

	// Output:
	// Created job: concat (one struct input and one struct output)
}
