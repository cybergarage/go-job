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
	"fmt"
	"testing"

	"github.com/cybergarage/go-job/job"
)

func TestArgumentsFrom(t *testing.T) {
	tests := []struct {
		v any
	}{
		{v: nil},
		{v: []any{1, 2, 3}},
		{v: []any{"a", "b", "c"}},
		{v: []string{"a", "b", "c"}},
		{v: "[1, 2, 3]"},
		{v: "[\"a\", \"b\", \"c\"]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("ArgumentsFrom %T", tt.v), func(t *testing.T) {
			_, err := job.NewArgumentsFrom(tt.v)
			if err != nil {
				t.Errorf("NewArgumentsFrom(%T) returned error: %v", tt.v, err)
			}
		})
	}
}
