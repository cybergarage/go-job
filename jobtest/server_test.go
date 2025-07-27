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

func ServerAPIsTest(t *testing.T, job job.Server) {
	t.Helper()

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
