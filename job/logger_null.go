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

type nullLogger struct{}

// NewNullLogger creates a new instance of null logger.
func NewNullLogger() Logger {
	return &nullLogger{}
}

// Infof logs an informational message.
func (l *nullLogger) Infof(job Job, format string, args ...any) {
	// No operation
}

// Errorf logs an error message.
func (l *nullLogger) Errorf(job Job, format string, args ...any) {
	// No operation
}

// Debugf logs a debug message.
func (l *nullLogger) Debugf(job Job, format string, args ...any) {
	// No operation
}

// Warnf logs a warning message.
func (l *nullLogger) Warnf(job Job, format string, args ...any) {
	// No operation
}

// Fatalf logs a fatal message and exits the program.
func (l *nullLogger) Fatalf(job Job, format string, args ...any) {
	// No operation
}
