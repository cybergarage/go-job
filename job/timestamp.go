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
	"time"
)

const (
	TimestampFormat = time.RFC3339Nano
)

// Timestamp represents a point in time.
type Timestamp time.Time

// NewTimestamp creates a new Timestamp from the current time.
func NewTimestamp() Timestamp {
	return Timestamp(time.Now())
}

// NewTimestampFromTime creates a new Timestamp from a time.Time value.
func NewTimestampFromTime(t time.Time) Timestamp {
	return Timestamp(t)
}

// NewTimestampFromString creates a new Timestamp from a string representation of time.
func NewTimestampFromString(s string) (Timestamp, error) {
	t, err := time.Parse(TimestampFormat, s)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp(t), nil
}

// Time returns the time.Time representation of the Timestamp.
func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

// Equal checks if two Timestamps are equal.
func (t Timestamp) Equal(other Timestamp) bool {
	return time.Time(t).Equal(time.Time(other))
}

// String returns the string representation of the Timestamp in RFC3339 format.
func (t Timestamp) String() string {
	return time.Time(t).Format(TimestampFormat)
}
