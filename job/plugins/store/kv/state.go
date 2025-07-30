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
	"github.com/cybergarage/go-job/job/encoding"
)

// NewInstanceStateKeyFrom creates a new key for a job instance state.
func NewInstanceStateKeyFrom(uuid job.UUID, suffixes ...string) Key {
	return newKeyFromUUID(instanceStatePrefix, uuid, suffixes...)
}

// NewInstanceStateListKey creates a new list key for job instance states.
func NewInstanceStateListKey() Key {
	return Key(instanceStatePrefix)
}

// NewObjectFromInstanceState creates a new Object from a job instance state.
func NewObjectFromInstanceState(state job.InstanceState, keySuffixes ...string) (Object, error) {
	data, err := encoding.MapToJSON(state.Map())
	if err != nil {
		return nil, err
	}
	return &object{
		key:   NewInstanceStateKeyFrom(state.UUID(), keySuffixes...),
		value: []byte(data),
	}, nil
}

// NewInstanceStateFromBytes creates a job instance state from a byte slice.
func NewInstanceStateFromBytes(b []byte) (job.InstanceState, error) {
	m, err := encoding.MapFromJSON(string(b))
	if err != nil {
		return nil, err
	}
	return job.NewInstanceStateFromMap(m)
}
