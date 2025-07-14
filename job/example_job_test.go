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

func ExampleNewJob() {
	NewJob(
		WithKind("sum"),
		WithExecutor(func(a, b int) int { return a + b }),
	)
}

func ExampleNewJob_struct() {
	type SumOpt struct {
		a int
		b int
	}

	NewJob(
		WithKind("sum"),
		WithExecutor(func(opt SumOpt) int { return opt.a + opt.b }),
	)
}
