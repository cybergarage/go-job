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
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/cmd/cli"
)

func ServerAPIsTest(t *testing.T, client job.Client, server job.Server) {
	t.Helper()

	var wg sync.WaitGroup

	err := server.Manager().Clear()
	if err != nil {
		t.Fatalf("failed to clear job manager: %v", err)
	}

	// Register a job with the server

	kind := "sum"

	resHandler := func(ji job.Instance, responses []any) {
		wg.Done()
	}

	errHandler := func(ji job.Instance, err error) error {
		ji.Errorf("Error: %v", err)
		return err
	}

	j, err := job.NewJob(
		job.WithKind(kind),
		job.WithExecutor(func(a, b int) int { return a + b }),
		job.WithCompleteProcessor(resHandler),
		job.WithTerminateProcessor(errHandler),
	)
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}
	err = server.Manager().RegisterJob(j)
	if err != nil {
		t.Fatalf("Failed to register job: %v", err)
	}

	// Open a client to the server

	err = client.Open()
	if err != nil {
		t.Fatalf("failed to open job client: %v", err)
	}

	defer func() {
		err := client.Close()
		if err != nil {
			t.Errorf("failed to close job client: %v", err)
		}
	}()

	// Verify the server version

	version, err := client.GetVersion()
	if err != nil {
		t.Fatalf("failed to get version: %v", err)
	}
	if version != job.Version {
		t.Errorf("expected version %s, got %s", job.Version, version)
	}

	// List registered jobs

	jobs, err := client.ListRegisteredJobs()
	if err != nil {
		t.Fatalf("failed to list registered jobs: %v", err)
	}
	if len(jobs) == 0 {
		t.Fatal("expected at least one registered job")
	}

	// Schedule a job

	wg.Add(1)

	instance, err := client.ScheduleJob(kind, 1, 2)
	if err != nil {
		t.Fatalf("failed to schedule job: %v", err)
	}
	if instance == nil {
		t.Fatal("expected job instance to be non-nil")
	}

	wg.Wait()

	// Lookup job instance

	instances, err := client.LookupInstances(
		job.NewQuery(
			job.WithQueryUUID(instance.UUID()),
		),
	)
	if err != nil {
		t.Fatalf("failed to lookup job instances: %v", err)
	}
	switch len(instances) {
	case 1:
		if instances[0].UUID() != instance.UUID() {
			t.Errorf("expected job instance UUID %s, got %s", instance.UUID(), instance.UUID())
		}
		if instances[0].State() != job.JobCompleted {
			t.Errorf("expected job instance (%s:%s) to be completed, got %s", instance.Kind(), instance.UUID(), instance.State())
		}
	default:
		t.Fatalf("expected exactly one job instance, got %d", len(instances))
	}
}

func TestServerAPIs(t *testing.T) {
	newCmdClient := func() job.Client {
		client := job.NewCliClient()
		client.SetCommandExecutor(func(name string, args ...string) ([]byte, error) {
			rootCmd := cli.GetRootCommand()
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetArgs(args)
			err := cli.Execute()
			if err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		})
		cli.SetClient(client)
		return client
	}

	t.Setenv("PATH", fmt.Sprintf("%s/bin:%s", os.Getenv("GOPATH"),
		os.Getenv("PATH")))

	clients := []job.Client{
		job.NewGrpcClient(),
		job.NewCliClient(),
		newCmdClient(),
	}

	serversConstructors := []func(opts ...any) (job.Server, error){
		job.NewServer,
	}

	servers := []job.Server{}
	for _, constructor := range serversConstructors {
		server, err := constructor()
		if err != nil {
			t.Fatalf("failed to create job server: %v", err)
		}
		servers = append(servers, server)
	}

	for _, server := range servers {
		t.Run(fmt.Sprintf("server(%s)", server.Manager().Store().Name()), func(t *testing.T) {
			// Start the server

			err := server.Start()
			if err != nil {
				t.Fatalf("failed to start job server: %v", err)
			}
			defer func() {
				err := server.Stop()
				if err != nil {
					t.Errorf("failed to stop job server: %v", err)
				}
			}()

			// Client tests

			for _, cli := range clients {
				cli.SetPort(server.GRPCPort())
				t.Run(fmt.Sprintf("client(%s)", cli.Name()), func(t *testing.T) {
					ServerAPIsTest(t, cli, server)
				})
			}

			// Prometheus metrics endpoint health check

			metricsURL := fmt.Sprintf("http://localhost:%d/metrics", server.PrometheusPort())
			resp, err := http.Get(metricsURL)
			if err != nil {
				t.Fatalf("failed to access Prometheus metrics endpoint: %v", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("unexpected status code from /metrics: %d", resp.StatusCode)
			}
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
		})
	}
}
