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

package cli

import (
	"fmt"

	"github.com/cybergarage/go-job/job"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listJobsCmd)
	listCmd.AddCommand(listInstancesCmd)
}

var listCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "list",
	Short: "List all resources",
	Long:  "List all resources in the specified category",
}

var listJobsCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "jobs",
	Short: "List registered jobs",
	Long:  "List all the registered jobs.",
	Run: func(cmd *cobra.Command, args []string) {
		jobs, err := GetClient().ListRegisteredJobs()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for n, job := range jobs {
			fmt.Printf("[%d] %s\n", n, job.String())
		}
	},
}

var listInstancesCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "instances",
	Short: "List job instances",
	Long:  "List all job instances.",
	Run: func(cmd *cobra.Command, args []string) {
		query := job.NewQuery()
		instances, err := GetClient().LookupInstances(query)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for n, instance := range instances {
			fmt.Printf("[%d] %s\n", n, instance.String())
		}
	},
}
