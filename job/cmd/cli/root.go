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

var gRPCHost string
var gRPCPort int

var rootCmd = &cobra.Command{ // nolint:exhaustruct
	Use:               "jobctl",
	Version:           job.Version,
	Short:             "Job Control CLI",
	Long:              "jobctl is a command-line interface for managing jobs in the go-job framework.",
	DisableAutoGenTag: true,
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}

func Execute() error {
	client := job.NewClient()
	client.SetHost(gRPCHost)
	client.SetPort(gRPCPort)

	if err := client.Open(); err != nil {
		return err
	}

	SetClient(client)

	defer func() {
		if err := client.Close(); err != nil {
			fmt.Fprintf(rootCmd.OutOrStderr(), "%s\n", err)
			return
		}
	}()

	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&gRPCHost, "host", "localhost", fmt.Sprintf("gRPC host or address for a %v instance", job.ProductName))
	rootCmd.PersistentFlags().IntVar(&gRPCPort, "port", job.DefaultGrpcPort, fmt.Sprintf("gRPC port number for a %v instance", job.ProductName))
}
