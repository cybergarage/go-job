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

package kv

import (
	"errors"
	"fmt"
)

var (
	ErrNotExist = errors.New("not exist")
)

// NewErrObjectNotExist returns a new error that the object is not exist.
func NewErrObjectNotExist(key Key) error {
	return fmt.Errorf("object (%s) is %w ", key.String(), ErrNotExist)
}
