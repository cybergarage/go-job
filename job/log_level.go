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
)

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

const (
	logErrorString = "ERROR"
	logInfoString  = "INFO"
	logWarnString  = "WARN"
)

// NewLogLevelFrom creates a new LogLevel from a specified value.
func NewLogLevelFrom(a any) (LogLevel, error) {
	switch v := a.(type) {
	case LogLevel:
		return v, nil
	case string:
		return NewLogLevelFromString(v)
	default:
		return 0, fmt.Errorf("invalid log level value: %v", a)
	}
}

// NewLogLevelFromString returns the LogLevel corresponding to the given string.
func NewLogLevelFromString(s string) (LogLevel, error) {
	switch s {
	case logInfoString:
		return LogInfo, nil
	case logErrorString:
		return LogError, nil
	case logWarnString:
		return LogWarn, nil
	}
	return 0, fmt.Errorf("unknown log level: %s", s)
}

// Contains checks if the LogLevel contains another LogLevel.
func (l LogLevel) Contains(other LogLevel) bool {
	return (l & other) != 0
}

// String returns the string representation of the Log.
func (l LogLevel) String() string {
	switch l {
	case LogInfo:
		return logInfoString
	case LogError:
		return logErrorString
	case LogWarn:
		return logWarnString
	default:
		return "UNKNOWN"
	}
}
