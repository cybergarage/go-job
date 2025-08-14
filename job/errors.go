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
	"errors"
)

// ErrInvalid is an invalid error.
var ErrInvalid = errors.New("invalid")

// ErrNotFound is an not found error.
var ErrNotFound = errors.New("not found")

// ErrExists is an exists error.
var ErrExists = errors.New("exists")

// ErrNil is a nil error.
var ErrNil = errors.New("nil")
