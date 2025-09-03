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

	"github.com/google/uuid"
)

// UUID is a type alias for uuid.UUID to represent job instance UUIDs.
type UUID = uuid.UUID

// NewUUID generates a new UUID for a job instance.
func NewUUID() UUID {
	tuuid, err := uuid.NewV7()
	if err != nil {
		return UUID(tuuid)
	}
	return UUID(uuid.New())
}

// NewUUIDFrom creates a new UUID from a given value.
func NewUUIDFrom(a any) (UUID, error) {
	switch v := a.(type) {
	case uuid.UUID:
		return UUID(v), nil
	case string:
		return NewUUIDFromString(v)
	default:
		return UUID{}, fmt.Errorf("invalid UUID value: %v", a)
	}
}

// NewUUIDFromString creates a UUID from a string representation.
func NewUUIDFromString(s string) (UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return UUID{}, err
	}
	return UUID(u), nil
}
