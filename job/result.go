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
)

// ResultSet represents the result of a job execution.
type ResultSet []any

// NewResultWith creates a new Result instance with the provided values.
func NewResultWith(values []any) ResultSet {
	return ResultSet(values)
}

// String returns a string representation of the Result.
func (r ResultSet) String() string {
	if len(r) == 0 {
		return "[]"
	}
	result := "["
	for i, v := range r {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%v", v)
	}
	result += "]"
	return result
}
