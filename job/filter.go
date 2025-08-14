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

// Filter is an interface that defines methods for filtering job instances.
type Filter interface {
	// Before returns the time before which job instances should be filtered.
	Before() (time.Time, bool)
	// After returns the time after which job instances should be filtered.
	After() (time.Time, bool)
	// IsUnset returns true if no filter criteria are configured.
	IsUnset() bool
	// Matches checks if the specified object matches the filter criteria.
	Matches(v any) bool
}

// FilterOption is a function that configures a job filter.
type FilterOption func(*filter)

type filter struct {
	before time.Time
	after  time.Time
}

// WithFilterBefore sets the time before which job instances should be filtered.
func WithFilterBefore(before time.Time) FilterOption {
	return func(f *filter) {
		f.before = before
	}
}

// WithFilterAfter sets the time after which job instances should be filtered.
func WithFilterAfter(after time.Time) FilterOption {
	return func(f *filter) {
		f.after = after
	}
}

// NewFilter creates a new instance of Filter with the given options.
func NewFilter(opts ...FilterOption) Filter {
	f := &filter{
		before: time.Time{},
		after:  time.Time{},
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Before returns the time before which job instances should be filtered.
func (f *filter) Before() (time.Time, bool) {
	if f.before.IsZero() {
		return time.Time{}, false
	}
	return f.before, true
}

// After returns the time after which job instances should be filtered.
func (f *filter) After() (time.Time, bool) {
	if f.after.IsZero() {
		return time.Time{}, false
	}
	return f.after, true
}

// IsUnset returns true if no filter criteria are configured.
func (f *filter) IsUnset() bool {
	if f == nil {
		return true
	}
	_, ok := f.Before()
	return !ok
}

// Matches checks if the specified object matches the filter criteria.
func (f *filter) Matches(v any) bool {
	if f.IsUnset() {
		return true
	}
	switch v := v.(type) {
	case Instance:
		if before, ok := f.Before(); ok && !v.CreatedAt().Before(before) {
			return false
		}
		if after, ok := f.After(); ok && !v.CreatedAt().After(after) {
			return false
		}
		return true
	case InstanceState:
		if before, ok := f.Before(); ok && !v.Timestamp().Before(before) {
			return false
		}
		if after, ok := f.After(); ok && !v.Timestamp().After(after) {
			return false
		}
		return true
	case Log:
		if before, ok := f.Before(); ok && !v.Timestamp().Before(before) {
			return false
		}
		if after, ok := f.After(); ok && !v.Timestamp().After(after) {
			return false
		}
		return true
	}
	return false
}
