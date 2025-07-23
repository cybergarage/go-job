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

// Map returns a map representation of the instance map.
func (im InstanceMap) Map() map[string]any {
	return map[string]any(im)
}
