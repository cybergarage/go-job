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

// NewInstanceKeyFromUUID creates a new key from a UUID string.
func NewInstanceKeyFrom(ji job.Instance, suffixes ...string) Key {
	return newKeyFrom(instancePrefix, suffixes...)
}

// NewInstanceListKey creates a new list key for a list of job instances.
func NewInstanceListKey() Key {
	return Key(instancePrefix)
}

// NewObjectFromInstance creates a new Object from a job instance.
func NewObjectFromInstance(ji job.Instance, suffixes ...string) (Object, error) {
	data, err := encoding.MapToJSON(ji.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to get JSON string from job instance: %w", err)
	}
	return &object{
		key:   NewInstanceKeyFrom(ji, suffixes...),
		value: []byte(data),
	}, nil
}

// NewInstanceFromBytes creates a job instance from a byte slice.
func NewInstanceFromBytes(b []byte, opts ...any) (job.Instance, error) {
	m, err := encoding.MapFromJSON(string(b))
	if err != nil {
		return nil, err
	}
	return job.NewInstanceFromMap(m, opts...)
}
