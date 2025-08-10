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

package kvutil

import (
	"testing"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
)

// LogObjects logs all key-value objects.
func LogObjects(t *testing.T, objects []kv.Object) {
	for n, obj := range objects {
		t.Logf("[%d] %s: %s\n", n, obj.Key().String(), string(obj.Bytes()))
	}
}
