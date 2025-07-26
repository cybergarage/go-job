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

// Kind is a type that represents the kind of a job.
type Kind = string

// NewKindFrom creates a new Kind from a specified value.
// It returns an error if the value is not a valid kind.
func NewKindFrom(a any) (Kind, error) {
	kind, ok := a.(string)
	if !ok {
		return "", fmt.Errorf("invalid kind value: %v", a)
	}
	return Kind(kind), nil
}
