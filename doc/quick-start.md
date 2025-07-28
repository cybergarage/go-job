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
		job.WithExecutor(func(a, b int) int { return a + b }),
		job.WithScheduleAt(time.Now()), // immediate scheduling is the default, so this option is redundant
		job.WithResponseHandler(func(ji job.Instance, res []any) {
			ji.Infof("Result: %v", res)
		}),
		job.WithErrorHandler(func(ji job.Instance, err error) error {
			ji.Errorf("Error executing job: %v", err)
			return err
		}),
	)
	mgr.RegisterJob(sumJob)

	// Schedule the registered job
	ji, _ := mgr.ScheduleRegisteredJob("sum", job.WithArguments(1, 2))

	// Start the job manager
	mgr.Start()

	// Wait for the job to complete
	mgr.StopWithWait()

	// Retrieve and print the job instance state history
	history, _ := mgr.LookupInstanceHistory(ji)
	for _, record := range history {
		fmt.Println(record.State())
	}

	// Retrieve and print the job instance logs
	logs, _ := mgr.LookupInstanceLogs(ji)
	for _, log := range logs {
		fmt.Println(log.Message())
	}
}
```
