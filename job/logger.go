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

// Logger is an interface that defines methods for logging messages.
type Logger interface {
	// Infof logs an informational message.
	Infof(job Job, format string, args ...any)
	// Errorf logs an error message.
	Errorf(job Job, format string, args ...any)
	// Debugf logs a debug message.
	Debugf(job Job, format string, args ...any)
	// Warnf logs a warning message.
	Warnf(job Job, format string, args ...any)
	// Fatalf logs a fatal message and exits the program.
	Fatalf(job Job, format string, args ...any)
}
