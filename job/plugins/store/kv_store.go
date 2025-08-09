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

package store

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kvutil"
)

type kvStore struct {
	kv.Store
}

func nowTimestampSuffix() string {
	return strconv.FormatInt(time.Now().UnixMicro(), 10)
}

// NewKvStoreWith creates a new key-value store instance.
func NewKvStoreWith(store kv.Store) job.Store {
	return &kvStore{
		Store: store,
	}
}

// EnqueueInstance stores a job instance in the store.
func (store *kvStore) EnqueueInstance(ctx context.Context, job job.Instance) error {
	obj, err := kv.NewObjectFromInstance(job)
	if err != nil {
		return err
	}
	return store.Set(ctx, obj)
}

// DequeueNextInstance retrieves and removes the highest priority job instance from the store. If no job instance is available, it returns nil.
func (store *kvStore) DequeueNextInstance(ctx context.Context) (job.Instance, error) {
	now := time.Now()

	rs, err := store.Scan(ctx, kv.NewInstanceListKey())
	if err != nil {
		return nil, err
	}

	var nextJob job.Instance
	for rs.Next() {
		obj, err := rs.Object()
		if err != nil {
			return nil, err
		}
		job, err := kv.NewInstanceFromBytes(
			obj.Bytes(),
		)
		if err != nil {
			return nil, err
		}
		scheduledAt := job.ScheduledAt()
		if scheduledAt.After(now) {
			continue
		}
		switch {
		case nextJob == nil:
			nextJob = job
		case job.Before(nextJob):
			nextJob = job
		}
	}

	if nextJob == nil {
		return nil, nil
	}

	key := kv.NewInstanceKeyFrom(nextJob)
	_, err = store.Remove(ctx, key)
	if err != nil {
		return nil, err
	}

	return nextJob, nil
}

// DequeueInstance removes a job instance from the store by its unique identifier.
func (store *kvStore) DequeueInstance(ctx context.Context, job job.Instance) error {
	key := kv.NewInstanceKeyFrom(job)
	_, err := store.Remove(ctx, key)
	return err
}

// ListInstances lists all job instances in the store.
func (store *kvStore) ListInstances(ctx context.Context) ([]job.Instance, error) {
	rs, err := store.Scan(ctx, kv.NewInstanceListKey())
	if err != nil {
		return nil, err
	}
	objs, err := kvutil.ReadAll(rs)
	if err != nil {
		return nil, err
	}
	jobs := []job.Instance{}
	for _, obj := range objs {
		job, err := kv.NewInstanceFromBytes(obj.Bytes())
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// ClearInstances clears all job instances in the store.
func (store *kvStore) ClearInstances(ctx context.Context) error {
	return store.Delete(ctx, kv.NewInstanceListKey())
}

// LogInstanceState adds a new state record for a job instance.
func (store *kvStore) LogInstanceState(ctx context.Context, state job.InstanceState) error {
	keySuffixes := []string{}
	if store.UniqueKeys() {
		keySuffixes = append(keySuffixes, nowTimestampSuffix())
	}
	obj, err := kv.NewObjectFromInstanceState(state, keySuffixes...)
	if err != nil {
		return err
	}
	return store.Set(ctx, obj)
}

// LookupInstanceHistory lists all state records for a job instance that match the specified query. The returned history is sorted by their timestamp.

func (store *kvStore) LookupInstanceHistory(ctx context.Context, query job.Query) (job.InstanceHistory, error) {
	rs, err := store.Scan(ctx, kv.NewInstanceStateListKey())
	if err != nil {
		return nil, err
	}
	objs, err := kvutil.ReadAll(rs)
	if err != nil {
		return nil, err
	}
	states := []job.InstanceState{}
	for _, obj := range objs {
		state, err := kv.NewInstanceStateFromBytes(obj.Bytes())
		if err != nil {
			return nil, err
		}
		if query.Matches(state) {
			states = append(states, state)
		}
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].Timestamp().Before(states[j].Timestamp())
	})
	return states, nil
}

// ClearInstanceHistory clears all state records for a job instance that match the specified filter.
func (store *kvStore) ClearInstanceHistory(ctx context.Context, filter job.Filter) error {
	if filter.IsUnset() {
		return store.Delete(ctx, kv.NewInstanceStateListKey())
	}
	rs, err := store.Scan(ctx, kv.NewInstanceStateListKey())
	if err != nil {
		return err
	}
	objs, err := kvutil.ReadAll(rs)
	if err != nil {
		return err
	}
	for _, obj := range objs {
		state, err := kv.NewInstanceStateFromBytes(obj.Bytes())
		if err != nil {
			return err
		}
		if !filter.Matches(state) {
			continue
		}
		err = store.Delete(ctx, kv.NewInstanceStateKeyFrom(state.UUID()))
		if err != nil {
			return err
		}
	}
	return nil
}

// Logf logs a formatted message at the specified log level.
func (store *kvStore) Logf(ctx context.Context, ji job.Instance, logLevel job.LogLevel, format string, args ...any) error {
	log := job.NewLog(
		job.WithLogKind(ji.Kind()),
		job.WithLogUUID(ji.UUID()),
		job.WithLogLevel(logLevel),
		job.WithLogMessage(fmt.Sprintf(format, args...)),
	)
	keySuffixes := []string{}
	if store.UniqueKeys() {
		keySuffixes = append(keySuffixes, nowTimestampSuffix())
	}
	obj, err := kv.NewObjectFromLog(log, keySuffixes...)
	if err != nil {
		return err
	}
	return store.Set(ctx, obj)
}

// Infof logs an informational message for a job instance.
func (store *kvStore) Infof(ctx context.Context, ji job.Instance, format string, args ...any) error {
	return store.Logf(ctx, ji, job.LogInfo, format, args...)
}

// Warnf logs a warning message for a job instance.
func (store *kvStore) Warnf(ctx context.Context, ji job.Instance, format string, args ...any) error {
	return store.Logf(ctx, ji, job.LogWarn, format, args...)
}

// Errorf logs an error message for a job instance.
func (store *kvStore) Errorf(ctx context.Context, ji job.Instance, format string, args ...any) error {
	return store.Logf(ctx, ji, job.LogError, format, args...)
}

// LookupInstanceLogs lists all log entries for a job instance. The returned logs are sorted by their timestamp.
func (store *kvStore) LookupInstanceLogs(ctx context.Context, query job.Query) ([]job.Log, error) {
	rs, err := store.Scan(ctx, kv.NewLogListKey())
	if err != nil {
		return nil, err
	}
	objs, err := kvutil.ReadAll(rs)
	if err != nil {
		return nil, err
	}
	logs := make([]job.Log, 0)
	for _, obj := range objs {
		log, err := kv.NewLogFromBytes(obj.Bytes())
		if err != nil {
			return nil, err
		}
		if query.Matches(log) {
			logs = append(logs, log)
		}
	}
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Timestamp().Before(logs[j].Timestamp())
	})
	return logs, nil
}

// ClearInstanceLogs clears all log entries for a job instance that match the specified filter.
func (store *kvStore) ClearInstanceLogs(ctx context.Context, filter job.Filter) error {
	if filter.IsUnset() {
		return store.Delete(ctx, kv.NewLogListKey())
	}

	rs, err := store.Scan(ctx, kv.NewLogListKey())
	if err != nil {
		return err
	}
	objs, err := kvutil.ReadAll(rs)
	if err != nil {
		return err
	}
	for _, obj := range objs {
		log, err := kv.NewLogFromBytes(obj.Bytes())
		if err != nil {
			return err
		}
		if !filter.Matches(log) {
			continue
		}
		err = store.Delete(ctx, kv.NewLogKeyFrom(log.UUID()))
		if err != nil {
			return err
		}
	}
	return nil
}

// Start starts the kv store.
func (store *kvStore) Start() error {
	return store.Store.Start()
}

// Stop stops the kv store.
func (store *kvStore) Stop() error {
	return store.Store.Stop()
}

// Clear clears all data in the kv store.
func (store *kvStore) Clear() error {
	return store.Store.Clear()
}
