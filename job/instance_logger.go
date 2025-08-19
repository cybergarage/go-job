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

// instanceLogger is an interface that defines methods for logging messages.
type instanceLogger interface {
	// Info logs an informational message.
	Info(msg string)
	// Warn logs a warning message.
	Warn(msg string)
	// Error logs an error message.
	Error(err error)
	// Infof logs an informational message with formatting.
	Infof(format string, args ...any)
	// Errorf logs an error message with formatting.
	Errorf(format string, args ...any)
	// Warnf logs a warning message with formatting.
	Warnf(format string, args ...any)
	// Debugf logs a debug message with formatting.
	Debugf(format string, args ...any)
}

// Info logs an informational message for the job instance.
func (ji *jobInstance) Info(msg string) {
	ji.history.Infof(ji, msg)
}

// Warn logs a warning message for the job instance.
func (ji *jobInstance) Warn(msg string) {
	ji.history.Warnf(ji, msg)
}

// Error logs an error message for the job instance.
func (ji *jobInstance) Error(err error) {
	ji.history.Errorf(ji, err.Error())
}

// Infof logs an informational message with formatting for the job instance.
func (ji *jobInstance) Infof(format string, args ...any) {
	ji.history.Infof(ji, format, args...)
}

// Warnf logs a warning message with formatting for the job instance.
func (ji *jobInstance) Warnf(format string, args ...any) {
	ji.history.Warnf(ji, format, args...)
}

// Errorf logs an error message with formatting for the job instance.
func (ji *jobInstance) Errorf(format string, args ...any) {
	ji.history.Errorf(ji, format, args...)
}

// Debugf logs a debug message with formatting for the job instance.
func (ji *jobInstance) Debugf(format string, args ...any) {
	ji.history.Debugf(ji, format, args...)
}
