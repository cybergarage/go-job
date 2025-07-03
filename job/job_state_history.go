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

// JobStateHistory keeps track of the state changes of a job.
type JobStateHistory struct {
	records []*StateRecord
}

// NewJobStateHistory creates a new job state history.
func NewJobStateHistory() *JobStateHistory {
	return &JobStateHistory{
		records: make([]*StateRecord, 0),
	}
}

// AppendStateRecord adds a new state record to the history with the current timestamp.
func (sh *JobStateHistory) AppendStateRecord(state JobState) {
	sh.records = append(sh.records, NewStateRecord(state))
}

// StateRecords returns the state records of the job.
func (sh *JobStateHistory) StateRecords() []*StateRecord {
	return sh.records
}

// LastStateRecord returns the most recent state record, or nil if there are no records.
func (sh *JobStateHistory) LastStateRecord() *StateRecord {
	if len(sh.records) == 0 {
		return nil
	}
	return sh.records[len(sh.records)-1]
}
