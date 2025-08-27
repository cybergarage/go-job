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
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cybergarage/go-safecast/safecast"
)

// Placeholder is a placeholder for the job arguments.
var Placeholder string = "?"

// Executor is a type that represents a function that executes a job.
// It can be any function type, allowing for flexible job execution.
type Executor any

// Execute calls the given function with the provided parameters and returns results as []any.
func Execute(fn any, args []any, opts ...any) (ResultSet, error) {
	fnObj := reflect.ValueOf(fn)
	fnType := fnObj.Type()
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("executor is not a function (%s)", fnType)
	}

	// inner variables

	var manager Manager
	var instance Instance
	var worker Worker
	var ctx context.Context

	for _, opt := range opts {
		switch v := opt.(type) {
		case Manager:
			manager = v
		case Instance:
			instance = v
		case Worker:
			worker = v
		case context.Context:
			ctx = v
		}
	}

	spArgs := map[reflect.Type]any{
		reflect.TypeOf((*context.Context)(nil)).Elem(): ctx,
		reflect.TypeOf((*Manager)(nil)).Elem():         manager,
		reflect.TypeOf((*Instance)(nil)).Elem():        instance,
		reflect.TypeOf((*Worker)(nil)).Elem():          worker,
	}

	// inner function for argument assignment

	argAssignable := func(arg any, fnArgType reflect.Type) (reflect.Value, bool) {
		argValue := reflect.ValueOf(arg)
		if argValue.IsValid() && argValue.Type().AssignableTo(fnArgType) {
			return argValue, true
		}
		return argValue, false
	}

	prepareArguments := func(fnType reflect.Type, args []any) ([]any, error) {
		if fnType.NumIn() == len(args) {
			return args, nil
		}
		prepArgs := []any{}
		argIdx := 0
		for n := range fnType.NumIn() {
			fnArgType := fnType.In(n)
			if argIdx < len(args) {
				arg := args[argIdx]
				if _, ok := argAssignable(arg, fnArgType); ok {
					prepArgs = append(prepArgs, arg)
					argIdx++
					continue
				}
			}

			switch fnArgType.Kind() {
			case reflect.Interface:
				v, ok := spArgs[fnArgType]
				if ok {
					prepArgs = append(prepArgs, v)
				}
			default:
				if argIdx < len(args) {
					arg := args[argIdx]
					prepArgs = append(prepArgs, arg)
					argIdx++
				}
			}
		}
		if len(prepArgs) != fnType.NumIn() {
			return nil, fmt.Errorf("argument count mismatch: want %d, got %d", fnType.NumIn(), len(prepArgs))
		}
		return prepArgs, nil
	}

	assignTo := func(arg any, fnArgType reflect.Type) (reflect.Value, bool) {
		assignMapTo := func(arg any, fnType reflect.Type) (reflect.Value, bool) {
			var argMap map[string]any

			switch arg := arg.(type) {
			case map[string]any:
				argMap = arg
			case string: /* JSON string */
				if err := json.Unmarshal([]byte(arg), &argMap); err != nil {
					return reflect.Value{}, false
				}
			case []byte: /* JSON bytes */
				if err := json.Unmarshal(arg, &argMap); err != nil {
					return reflect.Value{}, false
				}
			default:
				return reflect.Value{}, false
			}

			structValue := reflect.New(fnType).Elem()
			for mapKey, mapValue := range argMap {
				field := structValue.FieldByName(mapKey)
				if !field.IsValid() {
					return reflect.Value{}, false // Invalid field name
				}
				if !field.CanSet() {
					return reflect.Value{}, false
				}
				mapValue := reflect.ValueOf(mapValue)
				if mapValue.Type().ConvertibleTo(field.Type()) {
					field.Set(mapValue.Convert(field.Type()))
					continue
				}
				return reflect.Value{}, false
			}
			return structValue, true
		}

		argValue, ok := argAssignable(arg, fnArgType)
		if ok {
			return argValue, true
		}

		switch fnArgType.Kind() {
		case reflect.Struct:
			return assignMapTo(arg, fnArgType)
		case reflect.Ptr:
			switch fnArgType.Elem().Kind() {
			case reflect.Struct:
				structValue, ok := assignMapTo(arg, fnArgType.Elem())
				if !ok {
					return reflect.Value{}, false
				}
				ptrValue := reflect.New(fnArgType.Elem())
				ptrValue.Elem().Set(structValue)
				return ptrValue, true
			}
		case reflect.Array, reflect.Slice:
			return reflect.Value{}, false
		case reflect.Interface:
			v, ok := spArgs[fnArgType]
			if ok {
				return reflect.ValueOf(v), true
			}
			return reflect.Value{}, false
		default:
			fnVal := reflect.New(fnArgType).Interface()
			err := safecast.To(argValue.Interface(), fnVal)
			if err == nil {
				return reflect.ValueOf(fnVal).Elem(), true
			}
		}

		return reflect.Value{}, false
	}

	// execution (main function)

	args, err := prepareArguments(fnType, args)
	if err != nil {
		return nil, err
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
