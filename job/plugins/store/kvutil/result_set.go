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
	"github.com/cybergarage/go-job/job/plugins/store/kv"
)

// NewResultSetWithObjects creates a new ResultSet with the given objects.
func NewResultSetWithObjects(objs []kv.Object) kv.ResultSet {
	return &fullResultSet{
		objects: objs,
		cursor:  -1,
	}
}

type fullResultSet struct {
	objects []kv.Object
	cursor  int
}

// Next advances the cursor to the next object.
func (rs *fullResultSet) Next() bool {
	rs.cursor++
	return rs.cursor < len(rs.objects)
}

// Object returns the current object.
func (rs *fullResultSet) Object() (kv.Object, error) {
	if rs.cursor < 0 || len(rs.objects) <= rs.cursor {
		return nil, kv.ErrNotExist
	}
	return rs.objects[rs.cursor], nil
}
