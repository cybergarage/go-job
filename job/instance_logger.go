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

// InstanceLogger is an interface that defines methods for logging messages.
type InstanceLogger interface {
	// Info logs an informational message.
	Info(msg string)
	// Error logs an error message.
	Error(err error)
	// Debug logs a debug message.
	Debug(msg string)
	// Infof logs an informational message with formatting.
	Infof(format string, args ...any)
	// Errorf logs an error message with formatting.
	Errorf(format string, args ...any)
	// Debugf logs a debug message with formatting.
	Debugf(format string, args ...any)
}
