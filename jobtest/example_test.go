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
	"fmt"
	"time"

	"github.com/cybergarage/go-job/job"
)

func Example() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Register a job with a custom executor
	sumJob, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func(ji job.Instance, a, b int) int {
			ji.Infof("%s (%s): attempts %d", ji.UUID(), ji.Kind(), ji.Attempts())
			return a + b
		}),
		job.WithStateChangeProcessor(func(ji job.Instance, state job.JobState) {
			ji.Infof("State changed to: %v", state)
		}),
		job.WithCompleteProcessor(func(ji job.Instance, res []any) {
			ji.Infof("Result: %v", res)
		}),
		job.WithTerminateProcessor(func(ji job.Instance, err error) error {
			ji.Errorf("Error executing job: %v", err)
			return err
		}),
	)
	mgr.RegisterJob(sumJob)

	// Schedule the registered job
	ji, _ := mgr.ScheduleRegisteredJob("sum",
		job.WithArguments("?", 1, 2),   // "?" is a dummy placeholder for job.Instance
		job.WithScheduleAt(time.Now()), // immediate scheduling is the default, so this option is redundant
	)

	// Start the job manager
	mgr.Start()

	// Wait waits for all jobs to complete or terminate.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mgr.Wait(ctx)

	// Retrieve all queued and executed job instances
	query := job.NewQuery() // queries all job instances (any state)
	jis, _ := mgr.LookupInstances(query)
	for _, ji := range jis {
		fmt.Printf("Job Instance: %s, UUID: %s, State: %s\n", ji.Kind(), ji.UUID(), ji.State())
	}

	// Retrieve and print the job instance state history
	query = job.NewQuery(
		job.WithQueryInstance(ji), // filter by specific job instance
	)
	history, _ := mgr.LookupInstanceHistory(query)
	for _, record := range history {
		fmt.Println(record.State())
	}

	// Retrieve and print the job instance logs
	query = job.NewQuery(
		job.WithQueryInstance(ji), // filter by specific job instance
	)
	logs, _ := mgr.LookupInstanceLogs(query)
	for _, log := range logs {
		fmt.Println(log.Message())
	}

	// Stop the job manager
	mgr.Stop()
}
