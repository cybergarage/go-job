// Copyright (C) 2025 The go-job Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/cybergarage/go-job/job"
	"github.com/spf13/viper"
)

// Config represents a configuration interface.
type Config interface {
	// UseConfigFile uses the specified file as the configuration.
	UsedConfigFile() string
	// String returns a string representation of the configuration.
	String() string
}

type viperConfig struct {
}

// NewConfig creates a new configuration with the specified product name.
func NewConfig() Config {
	viper.SetConfigName(job.ProductName)
	viper.SetEnvPrefix(strings.ToUpper(job.ProductName))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	return &viperConfig{}
}

// UseConfigFile uses the specified file as the configuration.
func (conf *viperConfig) UsedConfigFile() string {
	return viper.ConfigFileUsed()
}

// String returns a string representation of the configuration.
func (conf *viperConfig) String() string {
	var s string
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		value := viper.Get(key)
		s += fmt.Sprintf("%s: %v\n", key, value)
	}
	return strings.TrimSuffix(s, "\n")
}
