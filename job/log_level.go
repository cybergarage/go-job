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

// LogLevel represents the level of log messages.
// It is used to categorize log messages for better organization and filtering.
type LogLevel int

const (
	// LogInfo represents informational log messages.
	LogInfo LogLevel = 1 << iota // 1
	// LogError represents error log messages.
	LogError // 2
	// LogWarn represents warning log messages.
	LogWarn // 4
	// LogAll represents all log levels combined.
	LogAll LogLevel = LogInfo | LogError | LogWarn // 7
)

// Contains checks if the LogLevel contains another LogLevel.
func (l LogLevel) Contains(other LogLevel) bool {
	return (l & other) != 0
}

// String returns the string representation of the Log.
func (l LogLevel) String() string {
	switch l {
	case LogInfo:
		return "INFO"
	case LogError:
		return "ERROR"
	case LogWarn:
		return "WARN"
	default:
		return "UNKNOWN"
	}
}
