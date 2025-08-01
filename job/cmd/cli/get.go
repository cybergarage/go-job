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

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getVersionCmd)
}

var getCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "get",
	Short: "Get the specified resource",
	Long:  "Get the specified resource in the specified category.",
}

var getVersionCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "version",
	Short: "Get version",
	Long:  "Get version string.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ver, err := GetClient().GetVersion()
		if err != nil {
			return err
		}
		fmt.Println(ver)
		return nil
	},
}
