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

// NewLogKeyFrom creates a new key for a job log.
func NewLogKeyFrom(uuid job.UUID) Key {
	return newKeyFromUUID(instanceLogPrefix, uuid)
}

// NewLogListKey creates a new list key for job instance logs.
func NewLogListKey() Key {
	return Key(instanceLogPrefix)
}

// NewObjectFromLog creates a new object from a job log entry.
func NewObjectFromLog(log job.Log) (Object, error) {
	data, err := encoding.MapToJSON(log.Map())
	if err != nil {
		return nil, err
	}
	return &object{
		key:   NewLogKeyFrom(log.UUID()),
		value: []byte(data),
	}, nil
}

// NewLogFromMap creates a new log entry from a map representation.
func NewInstanceLogFromBytes(b []byte) (job.Log, error) {
	m, err := encoding.MapFromJSON(string(b))
	if err != nil {
		return nil, err
	}
	return job.NewLogFromMap(m)
}
