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
	"fmt"
)

// InstanceMap is a map representation of a job instance.
type InstanceMap map[string]any

// NewInstanceMap creates a new empty instance map.
func NewInstanceMap() InstanceMap {
	return make(InstanceMap)
}

// NewInstanceMap creates a new instance with the provideed map.
func NewInstanceMapWith(m map[string]any) InstanceMap {
	return InstanceMap(m)
}

// Arguments returns the arguments from the instance map if they exist.
func (im InstanceMap) Arguments() (Arguments, bool) {
	if args, ok := im[argumentsKey]; ok {
		v, err := newArgumentsFrom(args)
		if err != nil {
			return nil, false
		}
		return v, true
	}
	return nil, false
}

// ResultSet returns the result set from the instance map if it exists.
func (im InstanceMap) ResultSet() (ResultSet, bool) {
	if rs, ok := im[resultSetKey]; ok {
		if resultSet, ok := rs.(ResultSet); ok {
			return resultSet, true
		}
	}
	return nil, false
}

// Error returns the error from the instance map if it exists.
func (im InstanceMap) Error() (error, bool) {
	if err, ok := im[errorKey]; ok {
		if errVal, ok := err.(error); ok {
			return errVal, true
		}
		if errStr, ok := err.(string); ok {
			return fmt.Errorf("%s", errStr), true
		}
	}
	return nil, false
}

// Map returns a map representation of the instance map.
func (im InstanceMap) Map() map[string]any {
	return map[string]any(im)
}
