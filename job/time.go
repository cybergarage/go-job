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
	"time"
)

const (
	TimestampFormat = time.RFC3339Nano
)

// Time represents a point in time.
type Time time.Time

// NewTime creates a new Timestamp from the current time.
func NewTime() Time {
	return Time(time.Now())
}

// NewTimeFrom creates a new Timestamp from a given value.
func NewTimeFrom(a any) (Time, error) {
	switch v := a.(type) {
	case Time:
		return v, nil
	case time.Time:
		return NewTimeFromTime(v), nil
	case string:
		return NewTimeFromString(v)
	default:
		return Time{}, fmt.Errorf("invalid timestamp value: %v", a)
	}
}

// NewTimeFromTime creates a new Timestamp from a time.Time value.
func NewTimeFromTime(t time.Time) Time {
	return Time(t)
}

// NewTimeFromString creates a new Timestamp from a string representation of time.
func NewTimeFromString(s string) (Time, error) {
	t, err := time.Parse(TimestampFormat, s)
	if err != nil {
		return Time{}, err
	}
	return Time(t), nil
}

// Time returns the time.Time representation of the Timestamp.
func (t Time) Time() time.Time {
	return time.Time(t)
}

// Equal checks if two Timestamps are equal.
func (t Time) Equal(other Time) bool {
	return time.Time(t).Equal(time.Time(other))
}

// String returns the string representation of the Timestamp in RFC3339 format.
func (t Time) String() string {
	return time.Time(t).Format(TimestampFormat)
}
