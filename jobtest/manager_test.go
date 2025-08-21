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

package jobtest

import (
	"context"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/plugins/store"
	"github.com/cybergarage/go-job/job/plugins/store/kv"
	"github.com/cybergarage/go-job/job/plugins/store/kvutil"
	"github.com/cybergarage/go-job/jobtest/plugins/store/kv/etcd"
	"github.com/cybergarage/go-job/jobtest/plugins/store/kv/redis"
	"github.com/cybergarage/go-job/jobtest/plugins/store/kv/valkey"
)

// nolint: maintidx
func ManagerJobScheduleTest(t *testing.T, mgr job.Manager) {
	t.Helper()

	startTimestamp := time.Now()

	if err := mgr.Start(); err != nil {
		t.Skipf("Failed to start job manager: %v", err)
		return
	}

	defer func() {
		if err := mgr.Stop(); err != nil {
			t.Errorf("Failed to stop job manager: %v", err)
		}
	}()

	if err := mgr.Clear(); err != nil {
		t.Errorf("Failed to clear job manager: %v", err)
		return
	}

	type sumOpt struct {
		A int
		B int
	}

	tests := []struct {
		kind string
		opts []any
		args []any
	}{
		{
			kind: "sum",
			opts: []any{
				job.WithExecutor(func(a, b int) int { return a + b }),
			},
			args: []any{1, 2},
		},
		{
			kind: "sum (string)",
			opts: []any{
				job.WithExecutor(func(a, b int) int { return a + b }),
			},
			args: []any{"1", "2"},
		},
		{
			kind: "sum (struct)",
			opts: []any{
				job.WithExecutor(func(opt sumOpt) int { return opt.A + opt.B }),
			},
			args: []any{sumOpt{1, 2}},
		},
		{
			kind: "sum (*struct)",
			opts: []any{
				job.WithExecutor(func(opt *sumOpt) int { return opt.A + opt.B }),
			},
			args: []any{&sumOpt{1, 2}},
		},
		{
			kind: "sum (instance)",
			opts: []any{
				job.WithExecutor(func(ji job.Instance, a, b int) int {
					ji.Debugf("Calculating %s: %d + %d", ji.Kind(), a, b)
					return a + b
				}),
			},
			args: []any{job.Placeholder, 1, 2},
		},
		{
			kind: "sum (manager+worker+instance)",
			opts: []any{
				job.WithExecutor(func(mgr job.Manager, ji job.Instance, w job.Worker, a, b int) int {
					ji.Debugf("NumWorkers: %d", mgr.NumWorkers())
					return a + b
				}),
			},
			args: []any{job.Placeholder, job.Placeholder, job.Placeholder, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			// Register a test job

			stateHandler := func(ji job.Instance, state job.JobState) {
				t.Logf("State changed to: %s (%s) %v", ji.Kind(), ji.UUID(), state)
				store := mgr.Store()
				switch store := store.(type) {
				case kv.Store:
					objs, err := store.Dump(context.Background())
					if err != nil {
						t.Errorf("Failed to dump store: %v", err)
					}
					kvutil.LogObjects(t, objs)
				}
			}

			processHandler := func(ji job.Instance, responses []any) {
				ji.Infof("%v", responses)
				ji.Errorf("%s: %v", job.LogError.String(), responses)
				ji.Warnf("%s: %v", job.LogWarn.String(), responses)
				wg.Done()
			}

			errorHandler := func(ji job.Instance, err error) error {
				ji.Errorf("Error: %v", err)
				t.Errorf("Error in job execution: %s (%s) %s", ji.Kind(), ji.UUID(), err)
				wg.Done()
				return err
			}

			opts := tt.opts
			opts = append(
				opts,
				job.WithKind(tt.kind),
				job.WithStateChangeProcessor(stateHandler),
				job.WithCompleteProcessor(processHandler),
				job.WithTerminateProcessor(errorHandler),
			)

			j, err := job.NewJob(opts...)
			if err != nil {
				t.Errorf("Failed to create job: %v", err)
				return
			}

			ji, err := mgr.RegisterJob(j)
			if err != nil {
				t.Errorf("Failed to register job: %v", err)
				return
			}
			if ji != nil {
				t.Errorf("Expected job instance to be nil, but got %v", ji)
				return
			}

			regJobs, err := mgr.ListJobs()
			if err != nil {
				t.Errorf("Failed to list registered jobs: %v", err)
				return
			}
			if len(regJobs) != 1 {
				t.Errorf("Expected exactly one registered job, but got %d", len(regJobs))
			}

			// Schedule the job with arguments

			ji, err = mgr.ScheduleRegisteredJob(
				tt.kind,
				job.WithScheduleAfter(0), // immediate scheduling
				job.WithArguments(tt.args...),
				job.WithBackoffDuration(0),
				job.WithTimeout(0),
			)
			if err != nil {
				t.Errorf("Failed to schedule job: %v", err)
				return
			}

			// Wait for the job to be processed

			wg.Wait()

			if err := mgr.Wait(); err != nil {
				t.Errorf("Failed to wait for jobs: %v", err)
				return
			}

			doneTimestamp := time.Now()

			// Lookup job instance (from history)

			instances, err := mgr.LookupInstances(job.NewQuery(job.WithQueryUUID(ji.UUID())))
			if err != nil {
				t.Errorf("Failed to lookup job instance: %v", err)
				return
			}
			switch len(instances) {
			case 1:
				if instances[0].UUID() != ji.UUID() {
					t.Errorf("Expected job instance UUID %s, but got %s", ji.UUID(), instances[0].UUID())
				}
				args := instances[0].Arguments()
				if len(args) != len(tt.args) {
					t.Errorf("Expected %d arguments, but got %d", len(tt.args), len(args))
				}
				_, err = instances[0].ResultSet()
				if err != nil {
					t.Errorf("Expected job instance to have a result set, but got error: %v", err)
				}
				attempt := instances[0].Attempts()
				if attempt != 1 {
					t.Errorf("Expected job instance to have 1 attempt, but got %d", attempt)
				}
			default:
				t.Errorf("Expected exactly one job instance, but got %d", len(instances))
			}

			// Check instance history

			query := job.NewQuery(
				job.WithQueryInstance(ji), // filter by specific job instance
			)
			history, err := mgr.LookupInstanceHistory(query)
			if err == nil {
				if len(history) == 0 {
					t.Errorf("Expected at least one history record for job instance")
				}

				lastState := history.LastState()
				if lastState == nil {
					t.Errorf("Expected last state to be non-nil, but it was nil")
					return
				}
				if lastState.State() != job.JobCompleted {
					t.Errorf("Expected last state to be %s, but got %s", job.JobCompleted, lastState.State())
				}
			} else {
				t.Errorf("Failed to retrieve instance history: %v", err)
				return
			}

			// Check history

			expectedStateOrders := []job.JobState{
				job.JobCreated,
				job.JobScheduled,
				job.JobProcessing,
				job.JobCompleted,
			}

			if len(history) != len(expectedStateOrders) {
				t.Errorf("Expected %d history records, but got %d", len(expectedStateOrders), len(history))
			}

			for i := range history {
				if history[i].State() != expectedStateOrders[i] {
					t.Errorf("Expected state %s at index %d, but got %s", expectedStateOrders[i], i, history[i].State())
				}
				if i == 0 {
					continue
				}
				pts := history[i-1].Timestamp()
				its := history[i].Timestamp()
				if pts.After(its) {
					t.Errorf("Record timestamps are not in non-decreasing order: record[%d]=%v, record[%d]=%v",
						i-1, pts, i, its)
				}
			}

			// Check logs

			expectedLogs := []string{
				"[3]",
			}

			query = job.NewQuery(
				job.WithQueryInstance(ji),          // filter by specific job instance
				job.WithQueryLogLevel(job.LogInfo), // filter logs by level
				job.WithQueryAfter(startTimestamp), // filter logs after the job started
				job.WithQueryBefore(doneTimestamp), // filter logs before the job finished
			)
			logs, err := mgr.LookupInstanceLogs(query)
			if err == nil {
				if len(logs) != len(expectedLogs) {
					t.Errorf("Expected %d logs, but got %d", len(expectedLogs), len(logs))
				}
				for i, log := range logs {
					if log.Message() != expectedLogs[i] {
						t.Errorf("Expected log %d to be %s, but got %s", i, expectedLogs[i], log.Message())
					}
				}
			} else {
				t.Errorf("Failed to retrieve instance logs: %v", err)
				return
			}

			// Unregister the job

			if err := mgr.UnregisterJob(tt.kind); err != nil {
				t.Errorf("Failed to unregister job: %v", err)
			}

			// Clean up

			query = job.NewQuery(
				job.WithQueryInstance(ji), // filter by specific job instance
			)

			if err := mgr.Clear(); err != nil {
				t.Errorf("Failed to clear job manager: %v", err)
			}

			history, err = mgr.LookupInstanceHistory(query)
			if err == nil {
				if len(history) != 0 {
					t.Errorf("Expected no history records after clearing, but got %d records", len(history))
				}
			} else {
				t.Errorf("Failed to retrieve instance history: %v", err)
			}

			logs, err = mgr.LookupInstanceLogs(query)
			if err == nil {
				if len(logs) != 0 {
					t.Errorf("Expected no logs after clearing, but got %d logs", len(logs))
				}
			} else {
				t.Errorf("Failed to retrieve instance logs: %v", err)
				return
			}
		})
	}
}

func ManagerJobCancelTest(t *testing.T, mgr job.Manager) {
	t.Helper()
}

func TestManager(t *testing.T) {
	tests := []func(t *testing.T, mgr job.Manager){
		ManagerJobScheduleTest,
		ManagerJobCancelTest,
	}

	stores := []job.Store{
		job.NewLocalStore(),
		store.NewMemdbStore(),
		store.NewKvStoreWith(valkey.NewStore()),
		store.NewKvStoreWith(etcd.NewStore()),
		store.NewKvStoreWith(redis.NewStore()),
	}

	for _, test := range tests {
		fnValue := reflect.ValueOf(test)
		testName := runtime.FuncForPC(fnValue.Pointer()).Name()
		t.Run(testName, func(t *testing.T) {
			for _, store := range stores {
				t.Run(store.Name(), func(t *testing.T) {
					mgr, err := job.NewManager(
						job.WithStore(store),
					)
					if err != nil {
						t.Errorf("Failed to create job manager: %v", err)
						return
					}
					test(t, mgr)
				})
			}
		})
	}
}
