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
	"github.com/google/uuid"
)

// NewInstancesFromQueue creates a list of job instances from the provided queue.
func NewInstancesFromQueue(queue Queue) ([]Instance, error) {
	list, err := queue.List()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// NewInstancesFromStore creates a list of job instances from the provided store.
func NewInstancesFromHistory(history InstanceHistory) ([]Instance, error) {
	jiOptsMap := make(map[uuid.UUID][]any)
	for _, state := range history {
		uuid := state.UUID()
		jiOpts, ok := jiOptsMap[uuid]
		if !ok {
			jiOpts = make([]any, 0)
			jiOpts = append(jiOpts, WithUUID(state.UUID()))
			jiOpts = append(jiOpts, WithKind(state.Kind()))
		}
		jiOpts = append(jiOpts, WithState(state.State()))
		stateMap := NewInstanceMapWith(state.Map())
		switch state.State() {
		case JobCreated:
			jiOpts = append(jiOpts, WithCreatedAt(state.Timestamp()))
			args, ok := stateMap.Arguments()
			if ok {
				jiOpts = append(jiOpts, WithArguments(args))
			}
		case JobScheduled:
			jiOpts = append(jiOpts, WithScheduleAt(state.Timestamp()))
		case JobProcessing:
			jiOpts = append(jiOpts, WithProcessingAt(state.Timestamp()))
		case JobCompleted:
			jiOpts = append(jiOpts, WithCompletedAt(state.Timestamp()))
			resultSet, ok := stateMap.ResultSet()
			if ok {
				jiOpts = append(jiOpts, WithResultSet(resultSet))
			}
		case JobTerminated:
			jiOpts = append(jiOpts, WithTerminatedAt(state.Timestamp()))
			err, ok := stateMap.Error()
			if ok {
				jiOpts = append(jiOpts, WithResultError(err))
			}
		case JobCancelled:
			jiOpts = append(jiOpts, WithCancelledAt(state.Timestamp()))
		case JobTimedOut:
			jiOpts = append(jiOpts, WithTimedOutAt(state.Timestamp()))
		}
		jiOptsMap[uuid] = jiOpts
	}
	jiList := make([]Instance, 0, len(jiOptsMap))
	for _, instanceOpts := range jiOptsMap {
		ji, err := NewInstance(instanceOpts...)
		if err != nil {
			return nil, err
		}
		jiList = append(jiList, ji)
	}
	return jiList, nil
}
