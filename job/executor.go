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

	"github.com/cybergarage/go-safecast/safecast"
)

// Executor is a type that represents a function that executes a job.
// It can be any function type, allowing for flexible job execution.
type Executor any

// Execute calls the given function with the provided parameters and returns results as []any.
func Execute(fn any, args ...any) (ResultSet, error) {
	fnObj := reflect.ValueOf(fn)
	fnType := fnObj.Type()
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("executor is not a function (%s)", fnType)
	}

	if fnType.NumIn() != len(args) {
		return nil, fmt.Errorf("argument count mismatch for function %s: want %d, got %d", fnType, fnType.NumIn(), len(args))
	}

	assignTo := func(arg any, fnType reflect.Type) (reflect.Value, bool) {
		argValue := reflect.ValueOf(arg)
		if argValue.Type().AssignableTo(fnType) {
			return argValue, true
		}

		fnVal := reflect.New(fnType).Interface()
		err := safecast.To(argValue.Interface(), fnVal)
		if err == nil {
			return reflect.ValueOf(fnVal).Elem(), true
		}

		return reflect.ValueOf(fnVal).Elem(), false
	}

	fnArgs := make([]reflect.Value, len(args))
	for i := range fnArgs {
		fnVal, ok := assignTo(args[i], fnType.In(i))
		if !ok {
			return nil, fmt.Errorf("argument[%d] type mismatch: want %v, got %v (%v)", i, fnType.In(i), reflect.TypeOf(args[i]), args[i])
		}
		fnArgs[i] = fnVal
	}

	reflectResults := fnObj.Call(fnArgs)
	results := make([]any, len(reflectResults))
	for i, r := range reflectResults {
		results[i] = r.Interface()
	}
	return results, nil
}
