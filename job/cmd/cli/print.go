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
	"github.com/cybergarage/go-job/job/encoding"
	"github.com/spf13/cobra"
)

func printJobs(cmd *cobra.Command, jobs []job.Job) error {
	cmd.Printf("[\n")
	for n, job := range jobs {
		json, err := encoding.MapToJSON(job.Map())
		if err != nil {
			return err
		}
		cmd.Printf("  %s", json)
		if n < len(jobs)-1 {
			cmd.Printf(",\n")
		} else {
			cmd.Printf("\n")
		}
	}
	cmd.Printf("]\n")
	return nil
}

func printInstance(cmd *cobra.Command, instance job.Instance) error {
	json, err := encoding.MapToJSON(instance.Map())
	if err != nil {
		return err
	}
	cmd.Println(json)
	return nil
}

func printInstances(cmd *cobra.Command, instances []job.Instance) error {
	cmd.Printf("[\n")
	for n, instance := range instances {
		json, err := encoding.MapToJSON(instance.Map())
		if err != nil {
			return err
		}
		cmd.Printf("  %s", json)
		if n < len(instances)-1 {
			cmd.Printf(",\n")
		} else {
			cmd.Printf("\n")
		}
	}
	cmd.Printf("]\n")
	return nil
}
