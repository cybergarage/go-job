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
	NewJob(
		WithKind("no args and no return"),
		WithDescription("A simple job that prints a message"),
		WithExecutor(func() { fmt.Println("Hello, World!") }),
	)
}

func ExampleNewJob_concat() {
	NewJob(
		WithKind("concat (two args and one return)"),
		WithDescription("Concatenates two strings"),
		WithExecutor(func(a, b string) string { return a + ", " + b }),
		WithArguments("hello", "world"),
	)
}

func ExampleNewJob_split() {
	NewJob(
		WithKind("split (one arg and two return)"),
		WithDescription("Splits a string into two parts"),
		WithExecutor(func(s string) (string, string) {
			parts := strings.Split(s, ",")
			return parts[0], parts[1]
		}),
		WithArguments("hello,world"),
	)
}

func ExampleNewJob_abs() {
	NewJob(
		WithKind("abs (one arg and one return)"),
		WithDescription("Returns the absolute value of an integer"),
		WithExecutor(func(a int) int { return int(math.Abs(float64(a))) }),
		WithArguments(-42),
	)
}

func ExampleNewJob_sum() {
	NewJob(
		WithKind("sum (two args and one return)"),
		WithDescription("Returns the sum of two integers"),
		WithExecutor(func(a, b int) int { return a + b }),
		WithArguments(1, 2),
	)
}

func ExampleNewJob_struct() {
	type SumOpt struct {
		a int
		b int
	}

	NewJob(
		WithKind("sum (struct arg and one return)"),
		WithDescription("Returns the sum of two integers from a struct"),
		WithExecutor(func(opt SumOpt) int { return opt.a + opt.b }),
		WithArguments(SumOpt{1, 2}),
	)
}
