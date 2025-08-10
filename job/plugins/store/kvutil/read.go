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
	"fmt"

	"github.com/cybergarage/go-job/job/plugins/store/kv"
)

// ReadAll reads all objects from the ResultSet and returns them as a slice.
func ReadAll(rs kv.ResultSet) ([]kv.Object, error) {
	var objects []kv.Object

	for rs.Next() {
		obj, err := rs.Object()
		if err != nil {
			return objects, fmt.Errorf("failed to read object: %w", err)
		}
		objects = append(objects, obj)
	}

	return objects, nil
}
