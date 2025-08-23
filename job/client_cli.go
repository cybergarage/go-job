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

package job

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const (
	jobctl = "jobctl"
)

type CLIClient interface {
	Client
	// SetCommandExecutor sets the command executor.
	SetCommandExecutor(executor CommandExecutor)
}

// CommandExecutor defines a function type for executing commands.
type CommandExecutor func(name string, args ...string) ([]byte, error)

// shell client implementation for client.
type cliClient struct {
	args     []string
	host     string
	port     int
	executor CommandExecutor
}

// NewCliClient returns a new cli client.
func NewCliClient(args ...string) CLIClient {
	client := &cliClient{
		args:     args,
		host:     "",
		port:     DefaultGRPCPort,
		executor: nil,
	}
	client.SetCommandExecutor(func(name string, args ...string) ([]byte, error) {
		cmd := exec.Command(name, args...)
		return cmd.CombinedOutput()
	})
	return client
}

// Name returns the name of the client.
func (cli *cliClient) Name() string {
	return jobctl
}

// SetPort sets a port number.
func (cli *cliClient) SetPort(port int) {
	cli.port = port
}

// SetHost sets a host name.
func (cli *cliClient) SetHost(host string) {
	cli.host = host
}

// SetCommandExecutor sets the command executor.
func (cli *cliClient) SetCommandExecutor(executor CommandExecutor) {
	cli.executor = executor
}

// Open opens a shell connection.
func (cli *cliClient) Open() error {
	_, err := exec.LookPath(jobctl)
	if err != nil {
		return fmt.Errorf("%s command not found in PATH: %w", jobctl, err)
	}
	return nil
}

// Close closes the shell connection.
func (cli *cliClient) Close() error {
	return nil
}

// Execute executes a command and returns the output.
func (cli *cliClient) Execute(name string, args ...string) ([]byte, error) {
	return cli.executor(name, args...)
}

func (cli *cliClient) trimOutput(out []byte) string {
	if len(out) > 0 && out[len(out)-1] == '\n' {
		return string(out[:len(out)-1])
	}
	return string(out)
}

// GetVersion retrieves the version of the job service.
func (cli *cliClient) GetVersion() (string, error) {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, cli.args...)
	cmdArgs = append(cmdArgs, "get", "version")
	out, err := cli.Execute(jobctl, cmdArgs...)
	if err != nil {
		return "", err
	}
	return cli.trimOutput(out), nil
}

// ScheduleJob schedules a job with the specified kind, priority, and arguments.
// The priority is lower for higher priority jobs, similar to Unix nice values.
func (cli *cliClient) ScheduleJob(kind string, args ...any) (Instance, error) {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, cli.args...)
	cmdArgs = append(cmdArgs, "schedule", kind)
	for _, arg := range args {
		cmdArgs = append(cmdArgs, fmt.Sprintf("%v", arg))
	}
	out, err := cli.Execute(jobctl, cmdArgs...)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(out, &m); err != nil {
		return nil, err
	}
	i, err := NewInstanceFromMap(m)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// ListRegisteredJobs lists all registered jobs.
func (cli *cliClient) ListRegisteredJobs() ([]Job, error) {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, cli.args...)
	cmdArgs = append(cmdArgs, "list", "jobs")
	out, err := cli.Execute(jobctl, cmdArgs...)
	if err != nil {
		return nil, err
	}
	var maps []map[string]any
	if err := json.Unmarshal(out, &maps); err != nil {
		return nil, err
	}
	jobs := make([]Job, len(maps))
	for n, m := range maps {
		j, err := newJobFromMap(m)
		if err != nil {
			return nil, err
		}
		jobs[n] = j
	}
	return jobs, nil
}

// LookupInstances looks up job instances based on the provided query.
func (cli *cliClient) LookupInstances(query Query) ([]Instance, error) {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, cli.args...)
	cmdArgs = append(cmdArgs, "list", "instances")
	out, err := cli.Execute(jobctl, cmdArgs...)
	if err != nil {
		return nil, err
	}
	var maps []map[string]any
	if err := json.Unmarshal(out, &maps); err != nil {
		return nil, err
	}
	instances := make([]Instance, len(maps))
	for n, m := range maps {
		i, err := NewInstanceFromMap(m)
		if err != nil {
			return nil, err
		}
		instances[n] = i
	}
	return instances, nil
}

// CancelInstances cancels job instances based on the provided query.
func (cli *cliClient) CancelInstances(query Query) ([]Instance, error) {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, cli.args...)
	cmdArgs = append(cmdArgs, "cancel", "instances")
	out, err := cli.Execute(jobctl, cmdArgs...)
	if err != nil {
		return nil, err
	}
	var maps []map[string]any
	if err := json.Unmarshal(out, &maps); err != nil {
		return nil, err
	}
	instances := make([]Instance, len(maps))
	for n, m := range maps {
		i, err := NewInstanceFromMap(m)
		if err != nil {
			return nil, err
		}
		instances[n] = i
	}
	return instances, nil
}
