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

package encoding

import (
	"encoding/json"
	"fmt"
)

// UnmarshalToMap converts any struct to map[string]any using JSON marshaling/unmarshaling.
func UnmarshalToMap(s any) map[string]any {
	errMap := func(err error) map[string]any {
		return map[string]any{"error": fmt.Sprintf("%T (%v)", s, err.Error())}
	}
	jsonData, err := json.Marshal(s)
	if err != nil {
		return errMap(err)
	}

	var m map[string]any
	err = json.Unmarshal(jsonData, &m)
	if err != nil {
		return errMap(err)
	}

	return m
}

// UnmarshalJSONToMap converts a JSON string to map[string]any.
func UnmarshalJSONToMap(jsonStr string) map[string]any {
	errMap := func(err error) map[string]any {
		return map[string]any{"error": fmt.Sprintf("JSON parse error: %v", err.Error())}
	}

	var m map[string]any
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		return errMap(err)
	}

	return m
}

// MergeMaps merges two maps into one, with values from m2 overwriting those in m1.
func MergeMaps(m1, m2 map[string]any) map[string]any {
	m := make(map[string]any)
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}
	return m
}

// MapToJSON converts a map[string]any to a JSON string.
func MapToJSON(m map[string]any) string {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("JSON marshal error: %v", err)
	}
	return string(jsonData)
}
