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
	"github.com/cybergarage/go-job/job/encoding"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scheduleCmd)
}

var scheduleCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "schedule",
	Short: "Schedule a job",
	Long:  "Schedule a job to run with the specified kind and arguments.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			cmd.Help()
			return errInvalidArguments(args)
		}

		kind := args[0]
		anyArgs := []any{}
		for _, arg := range args[1:] {
			anyArgs = append(anyArgs, arg)
		}

		job, err := GetClient().ScheduleJob(kind, anyArgs...)
		if err != nil {
			return err
		}

		json, err := encoding.MapToJSON(job.Map())
		if err != nil {
			return err
		}
		cmd.Println(json)

		return nil
	},
	Args:    cobra.MinimumNArgs(1), // Ensure at least one argument is provided
	Example: `job schedule kind args...`,
}
