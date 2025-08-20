# Quick Start Guide

### Installation

```sh
go get github.com/cybergarage/go-job
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/cybergarage/go-job/job"
)

func main() {
	// Create a job manager
	mgr, _ := job.NewManager()

	// Register a job with a custom executor
	sumJob, _ := job.NewJob(
		job.WithKind("sum"),
		job.WithExecutor(func(ji job.Instance, a, b int) int {
			ji.Debugf("sum(%d, %d)", a, b) // Log the input values using job.Instance method
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

	// Wait for the job to complete
	mgr.StopWithWait()

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
}
```
