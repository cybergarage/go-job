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
	"testing"

	"github.com/cybergarage/go-job/job"
)

func ServerAPIsTest(t *testing.T, server job.Server) {
	t.Helper()

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

	client := job.NewClient()

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
}

func TestServerAPIs(t *testing.T) {
	servers := []job.Server{}

	server, err := job.NewServer()
	if err != nil {
		t.Fatalf("failed to create job server: %v", err)
	}
	servers = append(servers, server)

	for _, srv := range servers {
		ServerAPIsTest(t, srv)
	}
}
