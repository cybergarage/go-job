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

package valkey

import (
	"github.com/cybergarage/go-job/job/plugins/store/kv"
)

// Memdb represents a Memdb instance.
type resultSet struct {
	nRead uint
}

func newResultSetWith() kv.ResultSet {
	return &resultSet{
		nRead: 0,
	}
}

// Next moves the cursor forward next object from its current position.
func (rs *resultSet) Next() bool {
	return false
}

// Object returns an object in the current position.
func (rs *resultSet) Object() (kv.Object, error) {
	return nil, nil
}
