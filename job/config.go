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

// Config is the interface for the job server configuration.
type Config interface {
	// SetGRPCPort sets the gRPC port for the job server.
	SetGRPCPort(port int)
	// GRPCPort returns the gRPC port for the job server.
	GRPCPort() int
	// SetPrometheusPort sets the Prometheus port for the job server.
	SetPrometheusPort(port int)
	// PrometheusPort returns the Prometheus port for the job server.
	PrometheusPort() int
}

type config struct {
	grpcPort       int
	prometheusPort int
}

// newConfig creates a new Config with default values.
func newConfig() *config {
	return &config{
		grpcPort:       DefaultGRPCPort,
		prometheusPort: DefaultPrometheusPort,
	}
}

// SetGRPCPort sets the gRPC port for the job server.
func (config *config) SetGRPCPort(port int) {
	config.grpcPort = port
}

// SetPrometheusPort sets the Prometheus port for the job server.
func (config *config) SetPrometheusPort(port int) {
	config.prometheusPort = port
}

// GRPCPort returns the gRPC port for the job server.
func (config *config) GRPCPort() int {
	return config.grpcPort
}

// PrometheusPort returns the Prometheus port for the job server.
func (config *config) PrometheusPort() int {
	return config.prometheusPort
}
