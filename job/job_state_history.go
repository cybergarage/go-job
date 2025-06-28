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
	records []*JobStateRecord
}

// NewJobStateHistory creates a new job state history.
func NewJobStateHistory() *JobStateHistory {
	return &JobStateHistory{
		records: make([]*JobStateRecord, 0),
	}
}

// AddRecord adds a new state record to the history.
func (sh *JobStateHistory) AddRecord(state JobState) {
	sh.records = append(sh.records, NewJobStateRecord(state))
}

// Records returns the state records of the job.
func (sh *JobStateHistory) Records() []*JobStateRecord {
	return sh.records
}

// LastRecord returns the most recent state record, or nil if there are no records.
func (sh *JobStateHistory) LastRecord() *JobStateRecord {
	if len(sh.records) == 0 {
		return nil
	}
	return sh.records[len(sh.records)-1]
}

// Clear removes all records from the history.
func (sh *JobStateHistory) Clear() {
	sh.records = make([]*JobStateRecord, 0)
}
