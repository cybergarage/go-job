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

type InstanceHistory interface {
	Records() []*InstanceRecord
}

// instanceHistory keeps track of the state changes of a job.
type instanceHistory struct {
	records []*InstanceRecord
}

// NewInstanceHistory creates a new job state history.
func NewInstanceHistory() *instanceHistory {
	return &instanceHistory{
		records: make([]*InstanceRecord, 0),
	}
}

// AppendRecord adds a new state record to the history with the current timestamp.
func (sh *instanceHistory) AppendRecord(state JobState) {
	sh.records = append(sh.records, NewInstanceRecord(state))
}

// Records returns the all instance records of the job.
func (sh *instanceHistory) Records() []*InstanceRecord {
	return sh.records
}

// LastRecord returns the most recent instance record, or nil if there are no records.
func (sh *instanceHistory) LastRecord() *InstanceRecord {
	if len(sh.records) == 0 {
		return nil
	}
	return sh.records[len(sh.records)-1]
}
