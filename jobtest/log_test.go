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

package jobtest

import (
	"testing"

	"github.com/cybergarage/go-job/job"
	"github.com/cybergarage/go-job/job/encoding"
)

func TestLogLevelFilter(t *testing.T) {
	tests := []struct {
		name     string
		level    job.LogLevel
		contains job.LogLevel
		expected bool
	}{
		{"Contains Info", job.LogAll, job.LogInfo, true},
		{"Contains Error", job.LogAll, job.LogError, true},
		{"Contains Warning", job.LogAll, job.LogWarn, true},
		{"Contains None", job.LogNone, job.LogAll, false},
		{"Does not contain Info", job.LogError, job.LogInfo, false},
		{"Does not contain Error", job.LogWarn, job.LogError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.Contains(tt.contains); got != tt.expected {
				t.Errorf("Contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLogLevelStrings(t *testing.T) {
	tests := []struct {
		level job.LogLevel
	}{
		{job.LogInfo},
		{job.LogError},
		{job.LogWarn},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			levelStr := tt.level.String()
			if levelStr == "" {
				t.Errorf("expected LogLevel %s to have a string representation", tt.level)
			}
			parsedLevel, err := job.NewLogLevelFromString(levelStr)
			if err != nil {
				t.Errorf("failed to parse LogLevel from string: %v", err)
			}
			if parsedLevel != tt.level {
				t.Errorf("expected LogLevel %s to equal parsed level %s", tt.level.String(), parsedLevel.String())
			}
			parsedLevelStr := parsedLevel.String()
			if parsedLevelStr == "" {
				t.Errorf("expected LogLevel %s to have a string representation", parsedLevel)
			}
		})
	}
}

func TestLogLJSON(t *testing.T) {
	tests := []struct {
		uuid    job.UUID
		level   job.LogLevel
		message string
	}{
		{uuid: job.NewUUID(), level: job.LogInfo, message: "This is an info message"},
		{uuid: job.NewUUID(), level: job.LogError, message: "This is an error message"},
		{uuid: job.NewUUID(), level: job.LogWarn, message: "This is a warning message"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			log := job.NewLog(
				job.WithLogUUID(tt.uuid),
				job.WithLogLevel(tt.level),
				job.WithLogMessage(tt.message),
			)

			jsonData, err := encoding.MapToJSON(log.Map())
			if err != nil {
				t.Errorf("failed to convert log to JSON: %v", err)
				return
			}

			jsonMap, err := encoding.MapFromJSON(string(jsonData))
			if err != nil {
				t.Errorf("failed to convert JSON to map: %v", err)
				return
			}

			parsedLog, err := job.NewLogFromMap(jsonMap)
			if err != nil {
				t.Errorf("failed to parse log from JSON: %v", err)
				return
			}

			if !log.Equal(parsedLog) {
				t.Errorf("expected log to equal parsed log, got %v, want %v", parsedLog, log)
			}
		})
	}
}
