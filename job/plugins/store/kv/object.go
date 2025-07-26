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
	"fmt"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/encoding"
)

// Object represents a key-value object.
type Object interface {
	// Key returns a key of the object.
	Key() Key
	// Value returns a value of the object.
	Value() []byte
}

type object struct {
	key   Key
	value []byte
}

// NewObjectFromInstance creates a new Object from a job instance.
func NewObjectFromInstance(ji job.Instance) (Object, error) {
	data, err := encoding.MapToJSON(ji.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to get JSON string from job instance: %w", err)
	}
	return &object{
		key:   NewInstanceKeyFrom(ji),
		value: []byte(data),
	}, nil
}

// Key returns the key of the object.
func (obj *object) Key() Key {
	return obj.key
}

// Value returns the value of the object.
func (obj *object) Value() []byte {
	return obj.value
}

// Map returns the object as a map.
func (obj *object) Map() (map[string]any, error) {
	return encoding.MapFromJSONString(string(obj.value))
}
