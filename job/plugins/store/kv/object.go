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
	"github.com/cybergarage/go-job/job"
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
func NewObjectFromInstance(instance job.Instance) Object {
	return &object{
		key:   Key(instance.UUID().String()),
		value: instance,
	}
}

// Key returns the key of the object.
func (obj *object) Key() Key {
	return obj.key
}

// Value returns the value of the object.
func (obj *object) Value() []byte {
	return obj.value
}
