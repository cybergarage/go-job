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

package jobtest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/cybergarage/go-job/job"
)

func TestExecutor(t *testing.T) {
	type concatString struct {
		a string
		b string
		s string
	}

	tests := []struct {
		fn       any
		params   []any
		expected []any
	}{
		{
			fn: func() {
				fmt.Println("Hello, World!")
			},
			params:   []any{},
			expected: nil,
		},
		{
			fn: func() int {
				return 42
			},
			params:   []any{},
			expected: []any{42},
		},
		{
			fn: func(s string) string {
				return "pong"
			},
			params:   []any{"ping"},
			expected: []any{"pong"},
		},
		{
			fn: func(s string) string {
				return s
			},
			params:   []any{"hello"},
			expected: []any{"hello"},
		},
		{
			fn: func(v1 int, v2 int) int {
				return v1 + v2
			},
			params:   []any{1, 2},
			expected: []any{3},
		},
		{
			fn: func(v1 int, v2 int) int {
				return v1 + v2
			},
			params:   []any{"1", "2"},
			expected: []any{3},
		},
		{
			fn: func(v1 int, v2 int) int {
				return v1 + v2
			},
			params:   []any{1, "2"},
			expected: []any{3},
		},
		{
			fn: func(a string, b string) string {
				return a + b
			},
			params:   []any{"foo", "bar"},
			expected: []any{"foobar"},
		},
		{
			fn: func(a string, b string) string {
				return a + b
			},
			params:   []any{123, 456},
			expected: []any{"123456"},
		},
		{
			fn: func(a float64, b float64) float64 {
				return a * b
			},
			params:   []any{2, 3.5},
			expected: []any{7.0},
		},
		{
			fn: func(a int) (int, int) {
				return a, a * 2
			},
			params:   []any{5},
			expected: []any{5, 10},
		},
		{
			fn: func(a bool) bool {
				return !a
			},
			params:   []any{false},
			expected: []any{true},
		},
		{
			fn:       func(a int) {},
			params:   []any{1},
			expected: []any{},
		},
		{
			fn: func(param *concatString) *concatString {
				param.s = param.a + " " + param.b
				return param
			},
			params:   []any{&concatString{"Hello", "world!", ""}},
			expected: []any{&concatString{"Hello", "world!", "Hello world!"}},
		},
	}

	for _, tt := range tests {
		paramsJSON, err := json.Marshal(tt.params)
		if err != nil {
			t.Fatalf("failed to marshal params: %v", err)
		}
		t.Run(fmt.Sprintf("%T(%s)", tt.fn, string(paramsJSON)), func(t *testing.T) {
			got, err := job.Execute(tt.fn, tt.params...)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d results, got %d", len(tt.expected), len(got))
				return
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.expected[i]) {
					t.Errorf("expected result[%d] to be %v, got %v", i, tt.expected[i], got[i])
				}
			}
		})
	}
}
