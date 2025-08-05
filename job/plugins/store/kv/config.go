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

package kv

// Config defines the configuration for the key-value store.
type Config interface {
	// UniqueKeys returns whether keys should be unique.
	UniqueKeys() bool
}

// ConfigOption defines a function that modifies the Config.
type ConfigOption func(*config)

type config struct {
	uniqueKeys bool
}

// WithUniqueKeys sets whether keys should be unique.
func WithUniqueKeys(unique bool) ConfigOption {
	return func(c *config) {
		c.uniqueKeys = unique
	}
}

// NewConfig creates a new Config with default values.
func NewConfig(opts ...ConfigOption) Config {
	c := &config{
		uniqueKeys: true, // default to unique keys
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// UniqueKeys returns whether keys should be unique.
func (c *config) UniqueKeys() bool {
	return c.uniqueKeys
}
