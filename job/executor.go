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
	"reflect"
)

// JobExecutor is a type that represents a function that executes a job.
type JobExecutor any

// Execute calls the given function with the provided parameters.
func Execute(fn any, params ...any) (result []reflect.Value, err error) {
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("executor is not a function")
	}
	if t.NumIn() != len(params) {
		return nil, fmt.Errorf("argument count mismatch")
	}
	in := make([]reflect.Value, len(params))
	for i, p := range params {
		paramType := t.In(i)
		val := reflect.ValueOf(p)
		if !val.Type().AssignableTo(paramType) {
			return nil, fmt.Errorf("argument[%d] type mismatch: want %v, got %v", i, paramType, val.Type())
		}
		in[i] = val
	}
	result = v.Call(in)
	return result, nil
}
