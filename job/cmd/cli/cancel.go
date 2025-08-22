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
	"github.com/cybergarage/go-job/job"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cancelCmd)
	cancelCmd.AddCommand(cancelInstancesCmd)
}

var cancelCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "cancel",
	Short: "cancel the specified resource",
	Long:  "cancel the specified resource in the specified query.",
}

var cancelInstancesCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "instances",
	Short: "Cancel job instances",
	Long:  "Cancel job instances by the specified query.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().StringP("kind", "k", "", "Kind of the instances to cancel")
		cmd.Flags().StringP("uuid", "u", "", "UUID of the instances to cancel")

		kind := cmd.Flags().Lookup("kind").Value.String()
		uuidStr := cmd.Flags().Lookup("uuid").Value.String()

		opts := []job.QueryOption{}
		if 0 < len(kind) {
			opts = append(opts, job.WithQueryKind(kind))
		}
		if 0 < len(uuidStr) {
			uuid, err := job.NewUUIDFrom(uuidStr)
			if err != nil {
				return err
			}
			opts = append(opts, job.WithQueryUUID(uuid))
		}
		instances, err := GetClient().CancelInstances(job.NewQuery(opts...))
		if err != nil {
			return err
		}
		return printInstances(cmd, instances)
	},
}
